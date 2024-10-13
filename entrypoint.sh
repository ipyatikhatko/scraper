#!/bin/bash

# Run the Go scraper once on container startup
echo "Running Go scraper on startup..."
/app/scraper

# Start cron daemon in the foreground with logging level 2
echo "Starting cron daemon..."
crond -f -l 2
