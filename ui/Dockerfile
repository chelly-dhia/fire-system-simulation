# Use Node.js official image
FROM node:18-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the Next.js application
RUN npm run build

# Install 'serve' to serve the built application
RUN npm install -g serve

# Expose port 3000 to the host
EXPOSE 3000

# Start the application using the production build
CMD ["serve", "-s", "out", "-l", "3000"]
