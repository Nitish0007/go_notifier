package interceptors

import (
	"log"
	"context"
	"strconv"
	
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
	"github.com/Nitish0007/go_notifier/internal/features/apiKey"
	accountv1 "github.com/Nitish0007/go_notifier/pkg/gen/account/v1"
)

// UnaryInterceptor is a gRPC unary interceptor that logs the request and response
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("UnaryInterceptor: ", info.FullMethod, "Request: ", req)
	
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %s: %v\n", info.FullMethod, r)
			_ = status.Errorf(codes.Internal, "Internal gRPC server error: %v\n", r)
		}
	}()

	resp, err := handler(ctx, req)
	if err != nil {
		log.Println("UnaryInterceptor error: ", err)
		return nil, err
	}
	return resp, err
}

func AuthUnaryInterceptor(db *gorm.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}
		
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		authKeys := md.Get("authorization")
		if len(authKeys) == 0 || authKeys[0] == "" {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
		authKey := authKeys[0]
		accountIDs := md.Get("account_id")
		if len(accountIDs) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing account_id")
		}
		accountID, err := strconv.ParseInt(accountIDs[0], 10, 64)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid account_id")
		}
		repo := apiKey.NewApiKeyRepository(db)
		keyRow, err := repo.FindByKeyAndAccountID(ctx, authKey, accountID)
		if err != nil || keyRow.Key == "" || keyRow.Key != authKey {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		ctx = context.WithValue(ctx, "account_id", accountID)
		return handler(ctx, req)
	}
}

func isPublicMethod(fullMethod string) bool {
	switch fullMethod {
		case accountv1.AccountService_Signup_FullMethodName, 
			accountv1.AccountService_Login_FullMethodName:
			return true
		default:
			return false
	}

}