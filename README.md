# Instruções

## Configuração

### Env
  - Criar arquivo .env na raiz do projeto com base no .env.example
    - RATE_LIMIT=Número máximo de requests por segundo por IP
    - DB_PORT=Porta para o banco de dados utilizado
    - DB_PASS=Senha para o banco de dados utilizado
    - DB_HOST=Host para o banco de dados utilizado
    - DB_NAME=Nome do banco de dados utilizado
    - BLOCK_TIME_MS=Tempo que um IP ou TOKEN ficam bloqueados após ultrapassarem o limite
### Tokens
- Na raiz do projeto existe um arquivo tokens.json onde são configurados limites para tokens especificos, é possível alterar o arquivo ou utilizar os tokens já cadastrados


## Execução

- Após a configuração é possível executar o projeto utilizando o seguinte comando `docker compose up --build`
- Para testar o limiter é possível utilizar `curl -H "API_KEY: Token2" http://localhost:8080`

## Informações

- É possível facilmente alterar o mecanismo de persistência apenas criando uma struct que implemente a interface de repositório `interfaces.Repository` e utiliza-la durante a criação do use case
- Foram realizados testes unitários para garantir a lógica do limiter
- Foram realizados testes de integração utilizando test containers para garantir a funcionalidade em situações de carga elevada