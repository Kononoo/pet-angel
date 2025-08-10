package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	avatv1 "pet-angel/api/avatar/v1"
	"pet-angel/internal/ai"
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	"pet-angel/internal/util"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// AvatarService 虚拟形象/道具/聊天 服务
// 负责：JWT 解析、参数与 proto 映射、错误日志

type AvatarService struct {
	avatv1.UnimplementedAvatarServiceServer
	uc        *biz.AvatarUsecase
	jwtSecret string
	logger    *log.Helper
}

// NewAvatarService 依赖注入构造器
func NewAvatarService(uc *biz.AvatarUsecase, authCfg *conf.Auth, l log.Logger) *AvatarService {
	secret := ""
	if authCfg != nil {
		secret = authCfg.JwtSecret
	}
	return &AvatarService{uc: uc, jwtSecret: secret, logger: log.NewHelper(l)}
}

// GetModels 获取可用模型
func (s *AvatarService) GetModels(ctx context.Context, in *avatv1.GetModelsRequest) (*avatv1.GetModelsReply, error) {
	list, err := s.uc.GetModels(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get models failed: %v", err)
		return nil, err
	}
	out := make([]*avatv1.PetModel, 0, len(list))
	for _, m := range list {
		out = append(out, &avatv1.PetModel{Id: m.ID, Name: m.Name, Path: m.Path, ModelType: m.ModelType, IsDefault: m.IsDefault, SortOrder: m.SortOrder})
	}
	return &avatv1.GetModelsReply{Models: out}, nil
}

// SetPetModel 设置当前模型
func (s *AvatarService) SetPetModel(ctx context.Context, in *avatv1.SetPetModelRequest) (*avatv1.SetPetModelReply, error) {
	userID, err := s.userIDFromCtx(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("set model: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.SetPetModel(ctx, userID, in.GetModelId()); err != nil {
		s.logger.WithContext(ctx).Errorf("set model: usecase error: %v", err)
		return nil, err
	}
	return &avatv1.SetPetModelReply{Success: true}, nil
}

// GetItems 获取道具列表
func (s *AvatarService) GetItems(ctx context.Context, in *avatv1.GetItemsRequest) (*avatv1.GetItemsReply, error) {
	list, err := s.uc.GetItems(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get items failed: %v", err)
		return nil, err
	}
	out := make([]*avatv1.Item, 0, len(list))
	for _, it := range list {
		out = append(out, &avatv1.Item{Id: it.ID, Name: it.Name, Description: it.Description, IconPath: it.IconPath, CoinCost: it.CoinCost})
	}
	return &avatv1.GetItemsReply{Items: out}, nil
}

// UseItem 使用道具（扣金币）
func (s *AvatarService) UseItem(ctx context.Context, in *avatv1.UseItemRequest) (*avatv1.UseItemReply, error) {
	userID, err := s.userIDFromCtx(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("use item: auth failed: %v", err)
		return nil, err
	}
	_, err = s.uc.UseItem(ctx, userID, in.GetItemId())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("use item: usecase error: %v", err)
		return nil, err
	}
	return &avatv1.UseItemReply{Success: true, Message: "ok"}, nil
}

// Chat 发送消息
func (s *AvatarService) Chat(ctx context.Context, in *avatv1.ChatRequest) (*avatv1.ChatReply, error) {
	userID, err := s.userIDFromCtx(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("chat: auth failed: %v", err)
		return nil, err
	}

	// 使用新的方法同时获取用户消息和AI回复
	userMsg, aiMsg, err := s.uc.GetChatWithAI(ctx, userID, in.GetContent())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("chat: usecase error: %v", err)
		return nil, err
	}

	// 构建回复
	reply := &avatv1.ChatReply{
		MessageId:   userMsg.ID,
		Content:     userMsg.Content,
		CreatedAt:   userMsg.CreatedAt.Format("2006-01-02 15:04:05"),
		AiMessageId: 0,
		AiContent:   "",
		AiCreatedAt: "",
	}

	// 如果有AI回复，添加到回复中
	if aiMsg != nil {
		reply.AiMessageId = aiMsg.ID
		reply.AiContent = aiMsg.Content
		reply.AiCreatedAt = aiMsg.CreatedAt.Format("2006-01-02 15:04:05")
	}

	return reply, nil
}

// ChatStream 流式聊天（GRPC版本，暂不支持流式）
func (s *AvatarService) ChatStream(ctx context.Context, in *avatv1.ChatRequest) (*avatv1.ChatReply, error) {
	// GRPC版本暂不支持真正的流式，返回错误提示使用HTTP版本
	return nil, biz.ErrNotImplemented
}

// ChatStreamHTTP 提供原生 HTTP SSE 处理器，便于在 server 中注册到指定路由
func (s *AvatarService) ChatStreamHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// 鉴权（demo模式放开校验）
		userID := int64(1) // 默认用户ID
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tok, err := jwtutil.FromAuthHeader(authHeader)
			if err == nil {
				if claims, err := jwtutil.Parse(s.jwtSecret, tok); err == nil {
					userID = claims.UserID
				}
			}
		}
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 先保存用户消息
		if _, err := s.uc.SaveUserMessage(r.Context(), userID, body.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		flusher, _ := w.(http.Flusher)

		// 获取AI客户端
		client := ai.Default()
		if client == nil {
			// 如果默认客户端为空，创建一个新的
			client = ai.NewClient(ai.Config{
				Model:       "deepseek-ai/DeepSeek-R1-0528-Qwen3-8B",
				BaseURL:     "https://api.siliconflow.cn/v1/",
				APIKey:      "sk-wucfvbppymimfcrmzrtbowbnpquyudkjbjpzahlavmlhddmq",
				MaxTokens:   16384,
				Temperature: 0.5,
			})
		}

		system := "你是一个治愈系的宠物数字伙伴，以第一人称'我'的口吻，温柔简短地回复。"
		var full string
		var hasContent bool

		// 调用AI流式接口
		_, streamErr := client.Stream(r.Context(), system, body.Content, func(delta string) error {
			if delta != "" {
				hasContent = true
				full += delta
				// 发送SSE数据
				_, _ = w.Write([]byte("data: " + delta + "\n\n"))
				if flusher != nil {
					flusher.Flush()
				}
			}
			return nil
		})

		// 处理错误
		if streamErr != nil {
			fmt.Printf("AI Stream Error: %v\n", streamErr)
			// 发送兜底回复
			fallbackReply := util.GetFallbackReply(body.Content)
			_, _ = w.Write([]byte("data: " + fallbackReply + "\n\n"))
			if flusher != nil {
				flusher.Flush()
			}
			full = fallbackReply
		} else if !hasContent {
			// 如果没有收到任何内容，发送兜底回复
			fallbackReply := util.GetFallbackReply(body.Content)
			_, _ = w.Write([]byte("data: " + fallbackReply + "\n\n"))
			if flusher != nil {
				flusher.Flush()
			}
			full = fallbackReply
		}

		// 保存AI回复到数据库
		if full != "" {
			if _, err := s.uc.SaveAIMessage(r.Context(), userID, full); err != nil {
				fmt.Printf("Save AI Message Error: %v\n", err)
			}
		}

		// 发送结束信号
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
		if flusher != nil {
			flusher.Flush()
		}
	}
}

// userIDFromCtx 从请求头解析 JWT 获取 user_id（demo模式放开校验）
func (s *AvatarService) userIDFromCtx(ctx context.Context) (int64, error) {
	// demo模式：如果没有token或token无效，返回默认用户ID
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return 1, nil // 默认用户ID
	}
	tok, err := jwtutil.FromAuthHeader(ts.RequestHeader().Get("Authorization"))
	if err != nil {
		return 1, nil // 默认用户ID
	}
	claims, err := jwtutil.Parse(s.jwtSecret, tok)
	if err != nil {
		return 1, nil // 默认用户ID
	}
	return claims.UserID, nil
}

// ensure time import used
var _ = time.Second
