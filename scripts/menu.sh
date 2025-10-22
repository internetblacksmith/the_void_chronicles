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
    echo "  1)  Run all tests"
    echo "  2)  Run tests with coverage report"
    echo "  3)  Run tests with verbose output"
    echo "  4)  Build the Go binary"
    echo "  5)  Run application locally (./run.sh)"
    echo "  6)  Lint and format code"
    echo "  7)  Run security scan"
    echo "  8)  Run all pre-commit checks"
    echo "  9)  Clean build artifacts"
    echo ""
    echo "🐳 Docker Commands:"
    echo "  10) Build Docker image locally"
    echo "  11) Run Docker container locally"
    echo ""
    echo "🚀 Deployment Commands (Kamal + Doppler):"
    echo "  12) 🔥 Deploy to production"
    echo "  13) Build and push image only"
    echo "  14) Stream production logs"
    echo "  15) Restart production containers"
    echo "  16) Rollback to previous version"
    echo "  17) Stop production containers"
    echo "  18) Open shell in production container"
    echo "  19) Show deployment status"
    echo "  20) Show production environment variables"
    echo "  21) Setup Kamal on new server"
    echo ""
    echo "  0)  Exit"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

execute_command() {
    case $1 in
        1)
            echo "▶️  Running tests..."
            make test
            ;;
        2)
            echo "▶️  Running tests with coverage..."
            make test-coverage
            ;;
        3)
            echo "▶️  Running tests with verbose output..."
            make test-verbose
            ;;
        4)
            echo "▶️  Building Go binary..."
            make build
            ;;
        5)
            echo "▶️  Running application locally..."
            make run
            ;;
        6)
            echo "▶️  Linting and formatting code..."
            make lint
            ;;
        7)
            echo "▶️  Running security scan..."
            make security-scan
            ;;
        8)
            echo "▶️  Running all pre-commit checks..."
            make pre-commit
            ;;
        9)
            echo "▶️  Cleaning build artifacts..."
            make clean
            ;;
        10)
            echo "▶️  Building Docker image..."
            make docker-build
            ;;
        11)
            echo "▶️  Running Docker container..."
            make docker-run
            ;;
        12)
            echo "🔥 Deploying to production..."
            echo ""
            read -p "⚠️  Are you sure you want to deploy to production? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy
            else
                echo "❌ Deployment cancelled."
            fi
            ;;
        13)
            echo "▶️  Building and pushing image..."
            make deploy-build
            ;;
        14)
            echo "▶️  Streaming production logs (Ctrl+C to exit)..."
            make deploy-logs
            ;;
        15)
            echo "▶️  Restarting production containers..."
            make deploy-restart
            ;;
        16)
            echo "▶️  Rolling back to previous version..."
            read -p "⚠️  Are you sure you want to rollback? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-rollback
            else
                echo "❌ Rollback cancelled."
            fi
            ;;
        17)
            echo "▶️  Stopping production containers..."
            read -p "⚠️  Are you sure you want to stop production? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-stop
            else
                echo "❌ Stop cancelled."
            fi
            ;;
        18)
            echo "▶️  Opening shell in production container..."
            make deploy-shell
            ;;
        19)
            echo "▶️  Showing deployment status..."
            make deploy-status
            ;;
        20)
            echo "▶️  Showing production environment variables..."
            make deploy-env
            ;;
        21)
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
        read -p "Select an option (0-21): " choice
        echo ""
        execute_command "$choice"
        echo ""
        if [ "$choice" != "5" ] && [ "$choice" != "11" ] && [ "$choice" != "14" ] && [ "$choice" != "18" ]; then
            read -p "Press Enter to continue..."
        fi
    done
}

main
