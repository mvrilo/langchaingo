package pinecone_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

func getValues(t *testing.T) (string, string, string, string) {
	t.Helper()

	pineconeAPIKey := os.Getenv("PINECONE_API_KEY")
	if pineconeAPIKey == "" {
		t.Skip("Must set PINECONE_API_KEY to run test")
	}

	environment := os.Getenv("PINECONE_ENVIRONMENT")
	if environment == "" {
		t.Skip("Must set PINECONE_ENVIRONMENT to run test")
	}

	indexName := os.Getenv("PINECONE_INDEX")
	if environment == "" {
		t.Skip("Must set PINECONE_INDEX to run test")
	}

	projectName := os.Getenv("PINECONE_PROJECT")
	if environment == "" {
		t.Skip("Must set PINECONE_INDEX to run test")
	}

	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	return environment, pineconeAPIKey, indexName, projectName
}

/* func TestPineconeStoreGRPC(t *testing.T) {
	t.Parallel()

	environment, apiKey, indexName, projectName := getValues(t)
	e, err := embeddings.NewOpenAI()
	require.NoError(t, err)

	storer, err := pinecone.New(
		context.Background(),
		pinecone.WithAPIKey(apiKey),
		pinecone.WithEnvironment(environment),
		pinecone.WithIndexName(indexName),
		pinecone.WithProjectName(projectName),
		pinecone.WithEmbedder(e),
		pinecone.WithNameSpace(uuid.New().String()),
		withGrpc(),
	)
	require.NoError(t, err)

	err = storer.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "yes"},
		{PageContent: "no"},
	})
	require.NoError(t, err)

	docs, err := storer.SimilaritySearch(context.Background(), "yeah", 1)
	require.NoError(t, err)
	require.Len(t, docs, 1)
	require.Equal(t, docs[0].PageContent, "yes")
} */

func TestPineconeStoreRest(t *testing.T) {
	t.Parallel()

	environment, apiKey, indexName, projectName := getValues(t)
	e, err := embeddings.NewOpenAI()
	require.NoError(t, err)

	storer, err := pinecone.New(
		context.Background(),
		pinecone.WithAPIKey(apiKey),
		pinecone.WithEnvironment(environment),
		pinecone.WithIndexName(indexName),
		pinecone.WithProjectName(projectName),
		pinecone.WithEmbedder(e),
		pinecone.WithNameSpace(uuid.New().String()),
	)
	require.NoError(t, err)

	err = storer.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "tokyo"},
		{PageContent: "potato"},
	})
	require.NoError(t, err)

	docs, err := storer.SimilaritySearch(context.Background(), "japan", 1)
	require.NoError(t, err)
	require.Len(t, docs, 1)
	require.Equal(t, docs[0].PageContent, "tokyo")
}

func TestPineconeStoreRestWithScoreThreshold(t *testing.T) {
	t.Parallel()

	environment, apiKey, indexName, projectName := getValues(t)
	e, err := embeddings.NewOpenAI()
	require.NoError(t, err)

	storer, err := pinecone.New(
		context.Background(),
		pinecone.WithAPIKey(apiKey),
		pinecone.WithEnvironment(environment),
		pinecone.WithIndexName(indexName),
		pinecone.WithProjectName(projectName),
		pinecone.WithEmbedder(e),
		pinecone.WithNameSpace(uuid.New().String()),
	)
	require.NoError(t, err)

	err = storer.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "Tokyo"},
		{PageContent: "Yokohama"},
		{PageContent: "Osaka"},
		{PageContent: "Nagoya"},
		{PageContent: "Sapporo"},
		{PageContent: "Fukuoka"},
		{PageContent: "Dublin"},
		{PageContent: "Paris"},
		{PageContent: "London "},
		{PageContent: "New York"},
	})
	require.NoError(t, err)

	// test with a score threshold of 0.8, expected 6 documents
	docs, err := storer.SimilaritySearch(context.Background(),
		"Which of these are cities in Japan", 10,
		vectorstores.WithScoreThreshold(0.8))
	require.NoError(t, err)
	require.Len(t, docs, 6)

	// test with a score threshold of 0, expected all 10 documents
	docs, err = storer.SimilaritySearch(context.Background(),
		"Which of these are cities in Japan", 10,
		vectorstores.WithScoreThreshold(0))
	require.NoError(t, err)
	require.Len(t, docs, 10)
}

func TestSimilaritySearchWithInvalidScoreThreshold(t *testing.T) {
	t.Parallel()

	environment, apiKey, indexName, projectName := getValues(t)
	e, err := embeddings.NewOpenAI()
	require.NoError(t, err)

	storer, err := pinecone.New(
		context.Background(),
		pinecone.WithAPIKey(apiKey),
		pinecone.WithEnvironment(environment),
		pinecone.WithIndexName(indexName),
		pinecone.WithProjectName(projectName),
		pinecone.WithEmbedder(e),
		pinecone.WithNameSpace(uuid.New().String()),
	)
	require.NoError(t, err)

	err = storer.AddDocuments(context.Background(), []schema.Document{
		{PageContent: "Tokyo"},
		{PageContent: "Yokohama"},
		{PageContent: "Osaka"},
		{PageContent: "Nagoya"},
		{PageContent: "Sapporo"},
		{PageContent: "Fukuoka"},
		{PageContent: "Dublin"},
		{PageContent: "Paris"},
		{PageContent: "London "},
		{PageContent: "New York"},
	})
	require.NoError(t, err)

	_, err = storer.SimilaritySearch(context.Background(),
		"Which of these are cities in Japan", 10,
		vectorstores.WithScoreThreshold(-0.8))
	require.Error(t, err)

	_, err = storer.SimilaritySearch(context.Background(),
		"Which of these are cities in Japan", 10,
		vectorstores.WithScoreThreshold(1.8))
	require.Error(t, err)
}

func TestPineconeAsRetriever(t *testing.T) {
	t.Parallel()

	environment, apiKey, indexName, projectName := getValues(t)
	e, err := embeddings.NewOpenAI()
	require.NoError(t, err)

	store, err := pinecone.New(
		context.Background(),
		pinecone.WithAPIKey(apiKey),
		pinecone.WithEnvironment(environment),
		pinecone.WithIndexName(indexName),
		pinecone.WithProjectName(projectName),
		pinecone.WithEmbedder(e),
	)
	require.NoError(t, err)

	id := uuid.New().String()

	err = store.AddDocuments(
		context.Background(),
		[]schema.Document{
			{PageContent: "The color of the house is blue."},
			{PageContent: "The color of the car is red."},
			{PageContent: "The color of the desk is orange."},
		},
		vectorstores.WithNameSpace(id),
	)
	require.NoError(t, err)

	llm, err := openai.New()
	require.NoError(t, err)

	result, err := chains.Run(
		context.TODO(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, 1, vectorstores.WithNameSpace(id)),
		),
		"What color is the desk?",
	)
	require.NoError(t, err)
	require.True(t, strings.Contains(result, "orange"), "expected orange in result")
}

func TestPineconeAsRetrieverWithScoreThreshold(t *testing.T) {
	t.Parallel()

	environment, apiKey, indexName, projectName := getValues(t)
	e, err := embeddings.NewOpenAI()
	require.NoError(t, err)

	store, err := pinecone.New(
		context.Background(),
		pinecone.WithAPIKey(apiKey),
		pinecone.WithEnvironment(environment),
		pinecone.WithIndexName(indexName),
		pinecone.WithProjectName(projectName),
		pinecone.WithEmbedder(e),
	)
	require.NoError(t, err)

	id := uuid.New().String()

	err = store.AddDocuments(
		context.Background(),
		[]schema.Document{
			{PageContent: "The color of the house is blue."},
			{PageContent: "The color of the car is red."},
			{PageContent: "The color of the desk is orange."},
			{PageContent: "The color of the lamp beside the desk is black."},
			{PageContent: "The color of the chair beside the desk is beige."},
		},
		vectorstores.WithNameSpace(id),
	)
	require.NoError(t, err)

	llm, err := openai.New()
	require.NoError(t, err)

	result, err := chains.Run(
		context.TODO(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, 5, vectorstores.WithNameSpace(
				id), vectorstores.WithScoreThreshold(0.8)),
		),
		"What colors is each piece of furniture next to the desk?",
	)
	require.NoError(t, err)

	require.Contains(t, result, "orange", "expected orange in result")
	require.Contains(t, result, "black", "expected black in result")
	require.Contains(t, result, "beige", "expected beige in result")
}
