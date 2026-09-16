# Mensago API

Serviço centralizado para criação, processamento e acompanhamento de campanhas de e-mail em lote.

O Mensago API permite que diferentes aplicações solicitem o envio de uma campanha por meio de uma API, sem precisar enviar cada mensagem individualmente ou manter regras próprias de fila, tentativas e integração com provedores de e-mail.

> Status do projeto: planejamento do MVP.

## Problema

Quando cada sistema implementa seu próprio mecanismo de envio, surgem problemas recorrentes:

- requisições lentas por causa do envio síncrono;
- duplicação de código entre aplicações;
- dificuldade para processar grandes listas de destinatários;
- ausência de retentativas e controle de falhas;
- dependência direta de um único provedor;
- pouca visibilidade sobre mensagens enviadas, pendentes ou rejeitadas.

O Mensago API concentra essas responsabilidades em um único serviço assíncrono e reutilizável.

## Objetivos

- Receber campanhas de diferentes aplicações clientes.
- Validar identidade e permissões com o Keycloak.
- Armazenar campanhas, destinatários e eventos de entrega.
- Processar os e-mails de forma assíncrona por meio de workers.
- Permitir a troca do provedor sem afetar as aplicações clientes.
- Aplicar retentativas, limites de envio e proteção contra duplicidade.
- Disponibilizar o estado e as métricas de cada campanha.

## Arquitetura

![Arquitetura inicial do Mensago API](docs/images/architecture.png)

### Componentes

| Componente | Responsabilidade |
| --- | --- |
| **App** | Sistema cliente que solicita o envio de uma campanha. |
| **User** | Usuário que cria, consulta ou administra campanhas. |
| **Campaign API** | Valida a requisição, cria a campanha e agenda seu processamento. |
| **Keycloak** | Autentica usuários e aplicações e fornece papéis e permissões. |
| **Database** | Armazena aplicações, campanhas, destinatários, tentativas e eventos. |
| **Worker** | Consome os envios pendentes, aplica limites e executa retentativas. |
| **Email Provider** | Realiza a entrega dos e-mails, por exemplo Amazon SES, Mailgun, Resend ou SMTP. |

## Fluxo de envio

1. Uma aplicação ou usuário obtém um token no Keycloak.
2. O cliente envia a campanha e os destinatários para a API.
3. A API valida os dados e registra a campanha no banco.
4. Os destinatários são divididos em tarefas de envio.
5. Um ou mais workers processam as tarefas de forma assíncrona.
6. O adaptador envia cada mensagem ao provedor configurado.
7. O resultado de cada tentativa é registrado no banco.
8. A aplicação consulta o andamento e as métricas da campanha.

## Escopo do MVP

- cadastro e identificação de aplicações clientes;
- criação de campanha imediata ou agendada;
- assunto, remetente, conteúdo HTML e conteúdo em texto;
- importação de destinatários;
- variáveis personalizadas por destinatário;
- processamento assíncrono em lote;
- controle de velocidade de envio;
- retentativas com intervalo progressivo;
- prevenção de envios duplicados;
- cancelamento de campanha ainda não concluída;
- consulta de progresso e falhas;
- integração inicial com um provedor de e-mail;
- recebimento de webhooks de entrega, rejeição e devolução.

## Estados principais

### Campanha

`draft` → `scheduled` → `processing` → `completed`

Uma campanha também poderá assumir os estados `paused`, `cancelled` ou `failed`.

### Mensagem

`pending` → `processing` → `sent` → `delivered`

Dependendo do retorno do provedor, uma mensagem também poderá ficar como `retrying`, `bounced`, `rejected` ou `failed`.

## API inicial

| Método | Rota | Descrição |
| --- | --- | --- |
| `POST` | `/v1/campaigns` | Cria uma campanha. |
| `GET` | `/v1/campaigns` | Lista campanhas da aplicação. |
| `GET` | `/v1/campaigns/{id}` | Retorna campanha, progresso e métricas. |
| `POST` | `/v1/campaigns/{id}/send` | Coloca uma campanha na fila de processamento. |
| `POST` | `/v1/campaigns/{id}/cancel` | Cancela uma campanha não concluída. |
| `GET` | `/v1/campaigns/{id}/messages` | Lista mensagens e seus estados. |
| `POST` | `/v1/webhooks/providers/{provider}` | Recebe eventos do provedor. |
| `GET` | `/health` | Verifica a disponibilidade do serviço. |

### Exemplo de criação

```http
POST /v1/campaigns
Authorization: Bearer <token>
Idempotency-Key: 3b757468-c61a-4b8d-bef1-8d014a0ae602
Content-Type: application/json
```

```json
{
  "name": "Boas-vindas aos alunos",
  "subject": "Bem-vindo, {{name}}!",
  "from": {
    "name": "WEBoot Cursos",
    "email": "contato@example.com"
  },
  "content": {
    "html": "<h1>Olá, {{name}}</h1><p>Sua inscrição foi confirmada.</p>",
    "text": "Olá, {{name}}. Sua inscrição foi confirmada."
  },
  "recipients": [
    {
      "email": "aluno@example.com",
      "variables": {
        "name": "João"
      }
    }
  ]
}
```

### Resposta esperada

```json
{
  "id": "01990e14-1e5f-7a78-a5a4-b920bd6228db",
  "status": "draft",
  "recipients": 1,
  "created_at": "2026-09-16T20:00:00Z"
}
```

## Stack sugerida

- **API:** Go
- **Roteamento HTTP:** Chi
- **Autenticação e autorização:** Keycloak / OpenID Connect
- **Banco de dados:** PostgreSQL
- **Fila:** NATS JetStream, RabbitMQ ou Redis Streams
- **Workers:** Go
- **Observabilidade:** OpenTelemetry, Prometheus e Grafana
- **Execução local:** Docker Compose

O mecanismo de fila deve ser decidido antes da implementação. O banco continuará sendo a fonte oficial do estado das campanhas e mensagens.

## Organização sugerida

```text
mensago-api/
├── cmd/
│   ├── api/
│   └── worker/
├── internal/
│   ├── auth/
│   ├── campaign/
│   ├── message/
│   ├── provider/
│   ├── queue/
│   └── webhook/
├── migrations/
├── docs/
│   └── images/
├── deployments/
├── .env.example
├── docker-compose.yml
├── go.mod
└── README.md
```

## Configuração prevista

```env
APP_ENV=development
HTTP_PORT=8080
DATABASE_URL=postgres://mensago:mensago@postgres:5432/mensago

KEYCLOAK_URL=http://keycloak:8080
KEYCLOAK_REALM=mensago
KEYCLOAK_CLIENT_ID=mensago-api

QUEUE_PROVIDER=nats
NATS_URL=nats://nats:4222

EMAIL_PROVIDER=mailgun
EMAIL_FROM_ADDRESS=contato@example.com
EMAIL_FROM_NAME=Mensago
```

Credenciais e chaves privadas não devem ser versionadas. Em produção, utilize secrets do ambiente de implantação.

## Requisitos importantes

- cada aplicação só pode consultar suas próprias campanhas;
- todas as rotas privadas devem validar token, aplicação e permissões;
- a criação e o disparo devem aceitar chave de idempotência;
- webhooks precisam ter assinatura validada;
- falhas temporárias devem gerar nova tentativa;
- falhas permanentes não devem ser reenviadas automaticamente;
- senhas, tokens e conteúdo sensível não devem aparecer nos logs;
- limites devem ser aplicados por aplicação e por provedor;
- descadastros e listas de bloqueio devem ser respeitados antes do envio.

## Roadmap

### Fase 1 — Fundação

- contratos da API e modelo de dados;
- autenticação com Keycloak;
- cadastro de aplicações;
- criação e consulta de campanhas.

### Fase 2 — Processamento

- fila e worker;
- integração com o primeiro provedor;
- idempotência, retentativas e limites;
- rastreamento individual das mensagens.

### Fase 3 — Entregabilidade

- webhooks do provedor;
- devoluções, rejeições e lista de bloqueio;
- métricas e observabilidade;
- configuração de SPF, DKIM e DMARC.

### Fase 4 — Evolução

- templates reutilizáveis;
- painel administrativo;
- múltiplos provedores com fallback;
- segmentação de contatos;
- suporte futuro a outros canais de comunicação.

## Uso responsável

O Mensago API deve ser utilizado apenas para comunicações autorizadas. As aplicações clientes precisam manter a origem do consentimento, oferecer uma forma clara de descadastro e observar a LGPD e as políticas dos provedores de e-mail.

## Licença

Este projeto está distribuído sob a licença MIT. Consulte o arquivo [LICENSE](LICENSE) para mais informações.
