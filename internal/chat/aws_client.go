package chat

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

func AWSConverse(question string, modelId string) (string, error) {
	ctx := context.Background()

	// Precedence: AWS_REGION > AWS_DEFAULT_REGION > 'us-east-1'
	var awsRegion string
	awsRegion = os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_DEFAULT_REGION")
	}
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		return "", err
	}

	client := bedrockruntime.NewFromConfig(cfg)
	resp, err := client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: aws.String(modelId),
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{Value: question},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	// Unwrap the response
	outputMsg := resp.Output.(*types.ConverseOutputMemberMessage)
	textBlock := outputMsg.Value.Content[0].(*types.ContentBlockMemberText)
	return textBlock.Value, nil
}
