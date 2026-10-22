FROM node:22.23.4-alpine AS deps
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@11.26.0 --activate
COPY code/frontend/package.json code/frontend/pnpm-lock.yaml code/frontend/pnpm-workspace.yaml* ./
ENV NODE_ENV=development
RUN pnpm install --frozen-lockfile

FROM node:22.23.4-alpine AS build
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@11.26.0 --activate
COPY --from=deps /app/node_modules ./node_modules
COPY code/frontend/ .

# Copy the changelog and documentation into the static folder so SvelteKit sees it during build
COPY CHANGELOG.md ./static/CHANGELOG.md
COPY docs ./static/docs
ARG VITE_API_BASE
ARG VITE_WS_BASE
ARG VITE_COMMIT_SHA
ARG VITE_REPO_NAME
ENV VITE_API_BASE=$VITE_API_BASE
ENV VITE_WS_BASE=$VITE_WS_BASE
ENV VITE_COMMIT_SHA=$VITE_COMMIT_SHA
ENV VITE_REPO_NAME=$VITE_REPO_NAME
RUN pnpm build

FROM nginxinc/nginx-unprivileged:1.29-alpine AS runtime
COPY --from=build /app/build /usr/share/nginx/html
# Fix nginx config path for root context
COPY code/frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
USER nginx
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
