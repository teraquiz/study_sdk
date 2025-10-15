# Teraquiz Study SDK

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

SDK standalone para consumir flashcards e questões diretamente do MongoDB do Teraquiz. Ideal para integração com outros serviços, ferramentas administrativas e aplicações que necessitam de acesso direto aos dados de estudo.

## Características

- Acesso direto e otimizado ao MongoDB
- Suporte a filtros avançados
- Gestão automática de conexões
- Zero dependências externas além do driver MongoDB
- Thread-safe e pronto para concorrência
- Totalmente tipado com Go

## Instalação

```bash
go get github.com/margarote/service_quiz_v2/pkg/study_sdk
```

## Uso

```go
package main

import (
    "context"
    "log"
    "time"

    sdk "github.com/teraquiz/study-service/pkg/study_sdk"
)

func main() {
    client, err := sdk.NewClient(sdk.Config{
        MongoURI:     "mongodb://localhost:27017",
        DatabaseName: "teraquiz",
        Timeout:      10 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close(context.Background())

    ctx := context.Background()

    flashcards, err := client.GetFlashcardsByCategory(ctx, "category-id")
    if err != nil {
        log.Fatal(err)
    }

    for _, fc := range flashcards {
        log.Printf("Front: %s\nBack: %s\n", fc.Front, fc.Back)
    }

    filter := sdk.FlashcardFilter{
        Language:   stringPtr("pt"),
        Difficulty: stringPtr("medium"),
        Enabled:    boolPtr(true),
    }

    filtered, err := client.ListFlashcards(ctx, filter)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found %d flashcards", len(filtered))
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }
```

## API

### Client

- `NewClient(config Config) (*Client, error)` - Cria cliente conectado ao MongoDB
- `Close(ctx context.Context) error` - Fecha conexão

### Métodos

- `GetFlashcardsByCategory(ctx, categoryID) ([]Flashcard, error)` - Busca por categoria
- `GetFlashcardByID(ctx, id) (*Flashcard, error)` - Busca por ID
- `ListFlashcards(ctx, filter) ([]Flashcard, error)` - Busca com filtros

### Filtros

```go
type FlashcardFilter struct {
    CategoryID *string
    Difficulty *string  // "easy", "medium", "hard"
    Language   *string  // ISO 639-1 (ex: "pt", "en")
    Verified   *bool
    Enabled    *bool
    Tags       []string
}
```
# study_sdk
