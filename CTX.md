# CTX
<!-- ctx-version: dev (built from source) -->
<!-- rules-version: 0.1.0 -->
<!-- schema-version: 0.1 -->
<!-- generated: 2025-08-20 -->

## ctx Overview

This project uses **ctx** - a universal tool wrapper that provides token awareness for AI agents. 

ctx wraps any CLI tool, shell command, or script to provide:
- Precise token counting (OpenAI, Anthropic, Gemini)
- Execution metadata and telemetry
- Safety controls (token limits, output limits)
- Structured JSON output for LLM consumption

**CRITICAL**: ALWAYS prefix commands with `ctx` to enable token-aware execution.

  INVOCATION:                                                                                                           
                                                                                                                        
    ctx [flags] run <command> [args...]    # Explicit subcommand                                                        
    ctx [flags] "<command with args>"      # Quoted command (simple cases)                                              
    ctx [flags] -- <command> [args...]     # POSIX standard separator                                                   
                                                                                                                        
  EXAMPLES:                                                                                                             
                                                                                                                        
    ctx --max-tokens 5000 run psql -c "SELECT * FROM users"                                                             
    ctx --no-tokens "echo Hello World"                                                                                  
    ctx --pretty run git status                                                                                         
    ctx --timeout 30s run long-running-script.sh                                                                        
                                                                                                                        
  ENVIRONMENT VARIABLES:                                                                                                
                                                                                                                        
    CLI flags take precedence over environment variables.                                                               
                                                                                                                        
    CTX_TOKEN_MODEL                                                                                                     
      Sets the default token counting provider (e.g., "anthropic", "openai", "gemini").                                 
                                                                                                                        
    CTX_NO_TOKENS                                                                                                       
      If "true", disables token counting for all commands.                                                              
                                                                                                                        
    CTX_PRETTY                                                                                                          
      If "true", outputs in pretty format instead of JSON.                                                              
                                                                                                                        
    CTX_MAX_TOKENS                                                                                                      
      Sets the maximum number of tokens allowed in output (e.g., "5000").                                               
                                                                                                                        
    CTX_MAX_OUTPUT_BYTES                                                                                                
      Sets the maximum number of bytes allowed in output (e.g., "1048576" for 1MB).                                     
                                                                                                                        
    CTX_MAX_LINES                                                                                                       
      Sets the maximum number of lines allowed in output (e.g., "1000").                                                
                                                                                                                        
    CTX_MAX_PIPELINE_STAGES                                                                                             
      Sets the maximum number of pipeline stages allowed (e.g., "5").                                                   
                                                                                                                        
    CTX_NO_HISTORY                                                                                                      
      If "true", disables command history recording.                                                                    
                                                                                                                        
    CTX_NO_TELEMETRY                                                                                                    
      If "true", disables OpenTelemetry tracing.                                                                        
                                                                                                                        
    CTX_PRIVATE                                                                                                         
      If "true", is equivalent to setting both CTX_NO_HISTORY and CTX_NO_TELEMETRY to "true".                           
                                                                                                                        
    CTX_TIMEOUT                                                                                                         
      Sets a default timeout for commands (e.g., "30s", "1m").                                                          
                                                                                                                        
    CTX_API_ENDPOINT                                                                                                    
      API endpoint for ctx Pro features (required for Pro features) (e.g., "https://api.ctx.click").                    
                                                                                                                        
    OTEL_EXPORTER_OTLP_ENDPOINT                                                                                         
      The OTLP endpoint for telemetry data (e.g., "http://localhost:4318").                                             
                                                                                                                        
                                                                                                                        
         
  USAGE  
         
    ctx <command> [command] [args...] [--flags]  
            
  COMMANDS  
            
    account                           View your ctx Pro account status
    completion [command]              Generate the autocompletion script for the specified shell
    config [command]                  Manage ctx configuration
    help [command]                    Help about any command
    login                             Authenticate with your ctx Pro account
    logout                            Log out of your ctx Pro account
    run [command] [args...]           Execute a command with ctx wrapping
    setup [command] [tool] [--flags]  Set up coding agents and assistants with ctx documentation
    telemetry [command]               Manage and view telemetry settings
    update [--flags]                  Update ctx to the latest version
    version [--flags]                 Show ctx version information
         
  FLAGS  
         
    -h --help                         Help for ctx
    --max-lines                       Maximum lines allowed in output (0 for no limit). Overrides CTX_MAX_LINES.
    --max-output-bytes                Maximum bytes allowed in output (0 for no limit). Overrides CTX_MAX_OUTPUT_BYTES.
    --max-pipeline-stages             Maximum pipeline stages allowed (0 for no limit). Overrides CTX_MAX_PIPELINE_STAGES.
    --max-tokens                      Maximum tokens allowed in output (0 for no limit). Overrides CTX_MAX_TOKENS.
    --no-history                      Disable saving command history. Overrides CTX_NO_HISTORY.
    --no-telemetry                    Disable OpenTelemetry tracing. Overrides CTX_NO_TELEMETRY.
    --no-tokens                       Disable token counting. Overrides CTX_NO_TOKENS.
    --output                          Output format ('json'). (json)
    --pretty                          Output in pretty format instead of JSON.
    --private                         Enable privacy mode (disables history and telemetry). Overrides CTX_PRIVATE.
    --stream                          Stream command output line by line for long-running tasks.
    --timeout                         Command execution timeout (e.g., '5s', '1m'). Overrides CTX_TIMEOUT. (0s)
    --token-model                     Token provider (anthropic, openai, gemini). Overrides CTX_TOKEN_MODEL.
    --version                         Show ctx version