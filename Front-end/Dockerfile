# Imagem base Node 22 (mais moderna e compatível com Angular 21)
FROM node:22-alpine

# Definir diretório de trabalho
WORKDIR /app

# Instalar Angular CLI globalmente
RUN npm install -g @angular/cli

# Copiar arquivos de configuração primeiro (se existirem)
COPY package*.json ./

# Instalar dependências
RUN npm install

# Copiar o restante do código
COPY . .

# Expor a porta padrão do Angular
EXPOSE 4200

# Comando para rodar o app em modo de desenvolvimento
# O host 0.0.0.0 é necessário para acesso externo ao container
CMD ["ng", "serve", "--host", "0.0.0.0", "--poll", "2000"]
