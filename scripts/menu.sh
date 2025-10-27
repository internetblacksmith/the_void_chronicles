#!/bin/bash

show_header() {
    clear
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  🚀 Void Chronicles - Development & Deployment Menu"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

show_menu() {
    echo "📦 Development Commands:"
    echo "  1)  Setup dev environment (install all dependencies)"
    echo "  2)  Run all tests"
    echo "  3)  Run tests with coverage report"
    echo "  4)  Run tests with verbose output"
    echo "  5)  Build the Go binary"
    echo "  6)  Run application locally (./run.sh)"
    echo "  7)  Lint and format code"
    echo "  8)  Run security scan"
    echo "  9)  Run all pre-commit checks"
    echo "  10) Clean build artifacts"
    echo "  11) Generate .kamal/secrets file"
    echo ""
    echo "🐳 Docker Commands:"
    echo "  12) Build Docker image locally"
    echo "  13) Run Docker container locally"
    echo ""
    echo "🚀 Deployment Commands (Kamal + Doppler):"
    echo "  14) 🔥 Deploy to production"
    echo "  15) Build and push image only"
    echo "  16) Stream production logs"
    echo "  17) Restart production containers"
    echo "  18) Rollback to previous version"
    echo "  19) Stop production containers"
    echo "  20) Open shell in production container"
    echo "  21) Show deployment status"
    echo "  22) Show production environment variables"
    echo "  23) Setup Kamal on new server"
    echo ""
    echo "  0)  Exit"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

execute_command() {
    case $1 in
        1)
            echo "▶️  Setting up development environment..."
            make setup
            ;;
        2)
            echo "▶️  Running tests..."
            make test
            ;;
        3)
            echo "▶️  Running tests with coverage..."
            make test-coverage
            ;;
        4)
            echo "▶️  Running tests with verbose output..."
            make test-verbose
            ;;
        5)
            echo "▶️  Building Go binary..."
            make build
            ;;
        6)
            echo "▶️  Running application locally..."
            make run
            ;;
        7)
            echo "▶️  Linting and formatting code..."
            make lint
            ;;
        8)
            echo "▶️  Running security scan..."
            make security-scan
            ;;
        9)
            echo "▶️  Running all pre-commit checks..."
            make pre-commit
            ;;
        10)
            echo "▶️  Cleaning build artifacts..."
            make clean
            ;;
        11)
            echo "▶️  Generating .kamal/secrets file..."
            make kamal-secrets-setup
            ;;
        12)
            echo "▶️  Building Docker image..."
            make docker-build
            ;;
        14)
            echo "🔥 Deploying to production..."
            echo ""
            read -p "⚠️  Are you sure you want to deploy to production? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy
            else
                echo "❌ Deployment cancelled."
            fi
            ;;
        15)
            echo "▶️  Building and pushing image..."
            make deploy-build
            ;;
        16)
            echo "▶️  Streaming production logs (Ctrl+C to exit)..."
            make deploy-logs
            ;;
        17)
            echo "▶️  Restarting production containers..."
            make deploy-restart
            ;;
        18)
            echo "▶️  Rolling back to previous version..."
            read -p "⚠️  Are you sure you want to rollback? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-rollback
            else
                echo "❌ Rollback cancelled."
            fi
            ;;
        19)
            echo "▶️  Stopping production containers..."
            read -p "⚠️  Are you sure you want to stop production? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-stop
            else
                echo "❌ Stop cancelled."
            fi
            ;;
        20)
            echo "▶️  Opening shell in production container..."
            make deploy-shell
            ;;
        21)
            echo "▶️  Showing deployment status..."
            make deploy-status
            ;;
        22)
            echo "▶️  Showing production environment variables..."
            make deploy-env
            ;;
        23)
            echo "▶️  Setting up Kamal on new server..."
            read -p "⚠️  Are you sure you want to setup a new server? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-setup
            else
                echo "❌ Setup cancelled."
            fi
            ;;
        0)
            echo "👋 Goodbye!"
            exit 0
            ;;
        *)
            echo "❌ Invalid option. Please try again."
            ;;
    esac
}

main() {
    while true; do
        show_header
        show_menu
        read -p "Select an option (0-23): " choice
        echo ""
        execute_command "$choice"
        echo ""
        if [ "$choice" != "6" ] && [ "$choice" != "13" ] && [ "$choice" != "16" ] && [ "$choice" != "20" ]; then
            read -p "Press Enter to continue..."
        fi
    done
}

main
