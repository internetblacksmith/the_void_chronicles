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
    echo "  1)  Setup dev environment (Go dependencies)           [make setup-dev]"
    echo "  2)  Setup deployment environment (Ruby/Kamal/Doppler) [make setup-deploy]"
    echo "  3)  Run all tests                                     [make test]"
    echo "  4)  Run tests with coverage report                    [make test-coverage]"
    echo "  5)  Run tests with verbose output                     [make test-verbose]"
    echo "  6)  Build the Go binary                               [make build]"
    echo "  7)  Run application locally (./run.sh)                [make run]"
    echo "  8)  Lint and format code                              [make lint]"
    echo "  9)  Run security scan                                 [make security-scan]"
    echo "  10) Run all pre-commit checks                         [make pre-commit]"
    echo "  11) Clean build artifacts                             [make clean]"
    echo "  12) Generate .kamal/secrets file                      [make kamal-secrets-setup]"
    echo ""
    echo "🐳 Docker Commands:"
    echo "  13) Build Docker image locally                        [make docker-build]"
    echo "  14) Run Docker container locally                      [make docker-run]"
    echo ""
    echo "🚀 Deployment Commands (Kamal + Doppler):"
    echo "  15) 🔥 Deploy to production                           [make deploy]"
    echo "  16) Build and push image only                         [make deploy-build]"
    echo "  17) Stream production logs                            [make deploy-logs]"
    echo "  18) Restart production containers                     [make deploy-restart]"
    echo "  19) Rollback to previous version                      [make deploy-rollback]"
    echo "  20) Stop production containers                        [make deploy-stop]"
    echo "  21) Open shell in production container                [make deploy-shell]"
    echo "  22) Show deployment status                            [make deploy-status]"
    echo "  23) Show production environment variables             [make deploy-env]"
    echo "  24) Setup Kamal on new server                         [make deploy-setup]"
    echo ""
    echo "  0)  Exit"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

execute_command() {
    case $1 in
        1)
            echo "▶️  Setting up development environment..."
            make setup-dev
            ;;
        2)
            echo "▶️  Setting up deployment environment..."
            make setup-deploy
            ;;
        3)
            echo "▶️  Running tests..."
            make test
            ;;
        4)
            echo "▶️  Running tests with coverage..."
            make test-coverage
            ;;
        5)
            echo "▶️  Running tests with verbose output..."
            make test-verbose
            ;;
        6)
            echo "▶️  Building Go binary..."
            make build
            ;;
        7)
            echo "▶️  Running application locally..."
            make run
            ;;
        8)
            echo "▶️  Linting and formatting code..."
            make lint
            ;;
        9)
            echo "▶️  Running security scan..."
            make security-scan
            ;;
        10)
            echo "▶️  Running all pre-commit checks..."
            make pre-commit
            ;;
        11)
            echo "▶️  Cleaning build artifacts..."
            make clean
            ;;
        12)
            echo "▶️  Generating .kamal/secrets file..."
            make kamal-secrets-setup
            ;;
        13)
            echo "▶️  Building Docker image..."
            make docker-build
            ;;
        15)
            echo "🔥 Deploying to production..."
            echo ""
            read -p "⚠️  Are you sure you want to deploy to production? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy
            else
                echo "❌ Deployment cancelled."
            fi
            ;;
        16)
            echo "▶️  Building and pushing image..."
            make deploy-build
            ;;
        17)
            echo "▶️  Streaming production logs (Ctrl+C to exit)..."
            make deploy-logs
            ;;
        18)
            echo "▶️  Restarting production containers..."
            make deploy-restart
            ;;
        19)
            echo "▶️  Rolling back to previous version..."
            read -p "⚠️  Are you sure you want to rollback? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-rollback
            else
                echo "❌ Rollback cancelled."
            fi
            ;;
        20)
            echo "▶️  Stopping production containers..."
            read -p "⚠️  Are you sure you want to stop production? (yes/no): " confirm
            if [ "$confirm" = "yes" ]; then
                make deploy-stop
            else
                echo "❌ Stop cancelled."
            fi
            ;;
        21)
            echo "▶️  Opening shell in production container..."
            make deploy-shell
            ;;
        22)
            echo "▶️  Showing deployment status..."
            make deploy-status
            ;;
        23)
            echo "▶️  Showing production environment variables..."
            make deploy-env
            ;;
        24)
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
        read -p "Select an option (0-24): " choice
        echo ""
        execute_command "$choice"
        echo ""
        if [ "$choice" != "7" ] && [ "$choice" != "14" ] && [ "$choice" != "17" ] && [ "$choice" != "21" ]; then
            read -p "Press Enter to continue..."
        fi
    done
}

main
