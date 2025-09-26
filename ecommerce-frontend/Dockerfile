# --- Step 1: Build Stage ---
# Uses the official Node.js 'slim' variant for a small and compatible build.
FROM node:22-slim AS builder

WORKDIR /app

# Copies the package.json and lockfile first to take advantage of Docker caching.
# 'npm install' will only be rerun if these files change.
COPY package*.json ./
RUN npm install

# Copies the rest of the frontend source code.
COPY . .

# Runs the application build to generate static files.
RUN npm run build

# --- Step 2: Production Stage ---
# Uses a super lightweight Nginx image to serve static content.
FROM nginx:stable-alpine AS production

# Copies the compiled files from the build stage to the Nginx web directory.
COPY --from=builder /app/dist /usr/share/nginx/html

# Applies a custom Nginx configuration to support Single-Page Application (SPA) routing.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Exposes port 80, the Nginx default.
EXPOSE 80

# Command to start the Nginx server in foreground mode.
CMD ["nginx", "-g", "daemon off;"]