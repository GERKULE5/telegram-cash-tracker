FROM oven/bun:1.4.2 AS dependencies
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile --production

FROM oven/bun:1.4.2
WORKDIR /app
ENV NODE_ENV=production
ENV WEB_PORT=3000
COPY --from=dependencies /app/node_modules ./node_modules
COPY package.json bunfig.toml ./
COPY src ./src
USER bun
EXPOSE 3000
CMD ["bun", "src/index.ts"]
