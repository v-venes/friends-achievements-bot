# Friends Achievements Bot 🎮🤖

Um bot de Discord desenvolvido em Go que envia diariamente as conquistas da Steam no seu servidor. O projeto utiliza RabbitMQ para comunicação assíncrona e Temporal para orquestração dos fluxos de extração de dados.

## 📌 Funcionalidades

- Comando no Discord para registrar o SteamID do usuário
- Rotina diária que consulta os jogos jogados recentemente na Steam
- Comparação automática de conquistas com o histórico do dia anterior
- Geração de imagem para cada nova conquista
- Envio automático da imagem no Discord via bot

---

## 🧱 Arquitetura

### 🔹 Discord Bot (Go)

- Responsável por escutar comandos no Discord (ex: `/adicionar-conta`)
- Escuta uma fila no RabbitMQ para receber eventos de novas conquistas
- Envia mensagens com imagens das conquistas para os canais apropriados

### 🔸 Queue Worker (Go)

- Escuta as demais filas do RabbitMQ
- Armazena SteamIDs dos usuários
- Roda uma rotina diária (orquestrada com **Temporal**) que:
  1. Busca os jogos jogados recentemente
  2. Coleta as conquistas atuais
  3. Compara com o snapshot do dia anterior
  4. Gera imagens das novas conquistas
  5. Publica evento na fila `new-achievement` do RabbitMQ

### 🐇 RabbitMQ

- Canal de comunicação entre o backend e o bot
- Fila: `new-achievement`

### ⏱️ Temporal

- Usado para orquestrar o fluxo diário de extração e verificação de conquistas
- Permite controle de retries, observabilidade, etc.

---

## 🚀 Como funciona o fluxo

```
1. [Usuário] usa comando /adicionar-conta no Discord
2. [Bot] envia SteamID para fila
3. [Worker] recebe a mensagem
  3.1 confere se o SteamID existe
  3.2 salva o player referente ao SteamID envia uma mensagem de feedback de sucesso caso exista
  3.3 envia uma mensagem de feedback de erro informando que o SteamID não existe
3. [Scheduler] inicia workflow diário (Temporal):
    ├─ Busca jogos recentes na Steam para os Players
    ├─ Coleta conquistas para cada jogo
    ├─ Compara com histórico
    ├─ Gera imagem
    └─ Publica evento na fila
4. [Bot] escuta fila `new-achievement`
5. [Bot] envia imagem da conquista no canal do Discord
```

---

## 🛠️ Tecnologias

- **Go**: linguagem principal
- **discordgo**: biblioteca para integração com Discord
- **RabbitMQ**: mensageria para comunicação assíncrona
- **Mongodb**: banco de dados geral
- **Temporal**: orquestração de workflows
- **Steam Web API**: coleta de dados dos jogos e conquistas

---

## 📦 Estrutura do repositório

```
.
├── discord-bot/         # Bot do Discord (escuta fila e envia mensagens)
├── queue_worker/             # Worker para processar demais mensagens
└── workflows/           # Workflows do Temporal
```

---

## 📄 Licença

Esse projeto é open source sob a licença [MIT](LICENSE).

---

## ✨ TODO

- [x] Criação do bot e comando para adicionar conta
- [x] Setup do RabbitMQ e conexão com o bot
- [x] Cadastro de SteamID via comando
  - [x] Verificar SteamID pela api
  - [x] Fazer integração com banco para guardar SteamID
  - [x] Enviar mensagem através do bot para informar a status
- [x] Workflow básico no Temporal
  - [x] Subir temporal no docker compose
  - [x] Conectar no temporal e registrar workflow
  - [ ] Extração dos últimos jogos para cada SteamID
  - [ ] Comparação de conquistas
  - [ ] Geração de imagem
- [ ] Integração completa bot ↔ backend
