# syntax=docker/dockerfile:1.7

FROM node:20-alpine AS builder
WORKDIR /src

RUN npm install -g pnpm@9.15.0

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

FROM node:20-alpine AS runtime
WORKDIR /app

RUN npm install -g pnpm@9.15.0

RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001

COPY --from=builder --chown=nextjs:nodejs /src/.output ./

USER nextjs

EXPOSE 3000

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST="0.0.0.0"

CMD ["node", "server/index.mjs"]
