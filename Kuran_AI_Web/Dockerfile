FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG VITE_API_BASE=https://kuranai-api.talhacan.com
ARG VITE_GOOGLE_CLIENT_ID=157832066999-3tgb2daavmrd1d1vd10ptrato1iu382t.apps.googleusercontent.com
ARG VITE_APPLE_CLIENT_ID=com.talhacan.kuranai.web
ARG VITE_APPLE_REDIRECT_URI=https://kuranai.talhacan.com/auth/apple/callback
ENV VITE_API_BASE=$VITE_API_BASE
ENV VITE_GOOGLE_CLIENT_ID=$VITE_GOOGLE_CLIENT_ID
ENV VITE_APPLE_CLIENT_ID=$VITE_APPLE_CLIENT_ID
ENV VITE_APPLE_REDIRECT_URI=$VITE_APPLE_REDIRECT_URI
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# SPA fallback: tüm route'ları index.html'e yönlendir
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
