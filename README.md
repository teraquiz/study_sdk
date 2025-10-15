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

    sdk "github.com/margarote/service_quiz_v2/pkg/study_sdk"
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

## Tipos

### Flashcard

```go
type Flashcard struct {
    ID         primitive.ObjectID // ID único do flashcard
    Language   string             // Código ISO 639-1 (pt, en, es)
    Front      string             // Frente do card (pergunta)
    Back       string             // Verso do card (resposta)
    Hint       string             // Dica opcional
    Difficulty string             // Dificuldade: easy, medium, hard
    Images     []FlashcardImage   // Imagens associadas
    Tags       []string           // Tags para categorização
    Enabled    bool               // Se o flashcard está ativo
    Verified   bool               // Se o flashcard foi verificado
    CreatedBy  string             // ID do criador
    VerifiedBy string             // ID do verificador
    CreatedAt  time.Time          // Data de criação
    UpdatedAt  time.Time          // Data de atualização
}

type FlashcardImage struct {
    URL     string // URL da imagem
    Caption string // Legenda opcional
}
```

## Casos de Uso

### Integração com Outros Serviços

Use o SDK quando você precisa:
- Consumir flashcards de outro microserviço
- Criar ferramentas administrativas
- Desenvolver scripts de migração/sincronização
- Implementar clientes que precisam acesso direto ao MongoDB

### Não Use o SDK Se

- Você está trabalhando **dentro** do `study-service` (use os repositories da camada de domínio)
- Você precisa de operações de escrita (SDK é read-only)
- Você quer seguir Clean Architecture no serviço principal

## Exemplos Avançados

### Busca com Múltiplos Filtros

```go
// Buscar flashcards verificados em português de dificuldade média
filter := sdk.FlashcardFilter{
    Language:   stringPtr("pt"),
    Difficulty: stringPtr("medium"),
    Verified:   boolPtr(true),
    Enabled:    boolPtr(true),
    Tags:       []string{"medicina", "anatomia"},
}

flashcards, err := client.ListFlashcards(ctx, filter)
if err != nil {
    log.Fatal(err)
}
```

### Processamento em Lote

```go
// Buscar por categoria e processar
flashcards, err := client.GetFlashcardsByCategory(ctx, categoryID)
if err != nil {
    log.Fatal(err)
}

for _, fc := range flashcards {
    // Processar cada flashcard
    fmt.Printf("ID: %s\n", fc.ID.Hex())
    fmt.Printf("Dificuldade: %s\n", fc.Difficulty)
    fmt.Printf("Tags: %v\n", fc.Tags)
}
```

## Configuração Avançada

### Timeout Customizado

```go
client, err := sdk.NewClient(sdk.Config{
    MongoURI:     "mongodb://localhost:27017",
    DatabaseName: "teraquiz",
    Timeout:      30 * time.Second, // Timeout de 30 segundos
})
```

### Usando com MongoDB Atlas

```go
client, err := sdk.NewClient(sdk.Config{
    MongoURI:     "mongodb+srv://user:pass@cluster.mongodb.net",
    DatabaseName: "teraquiz_prod",
    Timeout:      15 * time.Second,
})
```

## Requisitos

- Go 1.21 ou superior
- MongoDB 4.4 ou superior
- Acesso de leitura ao banco de dados Teraquiz

## Contribuindo

Contribuições são bem-vindas! Por favor:
1. Faça um fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## Suporte

Para problemas, dúvidas ou sugestões, abra uma [issue](https://github.com/margarote/service_quiz_v2/issues) no GitHub.
