FROM mcr.microsoft.com/playwright:v1.40.0-jammy

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy test files
COPY . .

# Install browsers
RUN npx playwright install --with-deps

# Set environment variables
ENV CI=true
ENV BASE_URL=http://test-app:8080

# Run tests
CMD ["npm", "test"]