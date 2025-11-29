#!/bin/bash
# Serve MkDocs documentation locally
# Usage: ./scripts/docs-serve.sh

set -e

echo "🚀 Starting MkDocs server..."
echo "📚 Documentation will be available at: http://127.0.0.1:8000"
echo "Press Ctrl+C to stop"
echo ""

mkdocs serve
