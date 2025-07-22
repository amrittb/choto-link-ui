# choto.link - URL Shortener

A modern, responsive URL shortening web application built with Go, Alpine.js, and Tailwind CSS. Features a beautiful Nepali-language interface with robust URL validation and user-friendly feedback.

## Features

- **Modern Responsive UI**: Clean design optimized for all devices
- **Nepali Language Interface**: Native Nepali text and user experience
- **URL Validation**: Robust input validation and sanitization
- **Loading & Copy Feedback**: Visual feedback for all user actions
- **Mobile Optimized**: Touch-friendly interface with proper scaling

## Quick Start

### Prerequisites
- Go 1.24.1 or higher
- Modern web browser

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/amrittb/choto-link-ui.git
   cd choto-link-ui
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run cmd/ui/main.go
   ```

4. **Open your browser**
   ```
   http://localhost:8080
   ```

## Technology Stack

### Backend
- **[Go](https://golang.org/)**: High-performance backend language
- **[Gin](https://gin-gonic.com/)**: Fast HTTP web framework
- **HTML Templates**: Server-side rendering

### Frontend
- **[Alpine.js](https://alpinejs.dev/)**: Lightweight reactive framework
- **[Tailwind CSS](https://tailwindcss.com/)**: Utility-first CSS framework
- **Vanilla JavaScript**: No heavy frameworks

### Development
- **Go Modules**: Dependency management
- **Hot Reload**: Manual restart required

## 🔧 Configuration

### Environment Variables
Currently, the application uses default configurations. Future versions may support:
- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: localhost)
- `DATABASE_URL`: Database connection string

## Deployment

### Local Development
```bash
go run cmd/ui/main.go
```

### Production Build
```bash
go build -o choto-link cmd/ui/main.go
./choto-link
```

### Docker (Future)
```dockerfile
# Dockerfile example for future use
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/ui/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/internal ./internal
CMD ["./main"]
```

## Testing

### Manual Testing
1. **Valid URLs**: Test with various valid URL formats
2. **Invalid URLs**: Test error handling with invalid inputs
3. **Mobile**: Test responsive design on different devices
4. **Copy Function**: Verify copy-to-clipboard functionality
5. **Loading States**: Check loading spinner and disabled states

### Future Automated Testing
- Unit tests for URL validation
- Integration tests for API endpoints
- E2E tests for user workflows

**Made with ❤️ for the Nepali community**

*choto.link - आफ्नो लिङ्क छोटो बनाउनुहोस्, सजिलै साझा गर्नुहोस्!*