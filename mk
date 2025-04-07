#!/bin/bash

# Function to fetch and prepare terraform providers for user-wide discovery using GitHub CLI
function download_providers() {
  # Array of plugin names
  local plugins=("$@")
  
  # Loop through each plugin
  for plugin in "${plugins[@]}"; do
    echo "Processing plugin: $plugin"
    
    # Define the GitHub repository for the plugin
    local repo="patrikkj/terraform-provider-${plugin}"

    # Use `gh release download` to download the release zip file for the plugin
    echo "Downloading $plugin..."
    gh release download --repo "$repo" --pattern "terraform-provider-${plugin}_0.0.1_darwin_arm64.zip" --clobber --dir .
    
    # Check if the zip file was downloaded
    if [ ! -f "terraform-provider-${plugin}_0.0.1_darwin_arm64.zip" ]; then
      echo "Error: Download failed for $plugin"
      exit 1
    fi

    # Extract the zip file
    echo "Extracting $plugin..."
    unzip -o "terraform-provider-${plugin}_0.0.1_darwin_arm64.zip"
    
    # Create the necessary directory in the Terraform plugin location
    local plugin_dir="$HOME/.terraform.d/plugins/patrikkj/${plugin}/0.0.1/darwin_amd64"
    mkdir -p "$plugin_dir"
    
    # Move the extracted binary to the appropriate directory
    mv "terraform-provider-${plugin}" "$plugin_dir/"
    
    # Clean up the zip file
    rm "terraform-provider-${plugin}_0.0.1_darwin_arm64.zip"
    
    echo "$plugin prepared for Terraform discovery."
  done
}

function help() {
    echo "Available commands:"
    perl -ne 'print "  $1\n" if /^function ([^_]\w*)/' "$0" | sort
}

# Handle command line arguments
if [ $# -eq 0 ]; then
    help
else
    # Execute the function with the provided arguments (first arg is fn name)
    $1 "${@:2}"
fi
