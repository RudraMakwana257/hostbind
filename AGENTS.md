> This repo uses HostBind. 
> Never choose or assume dev ports (e.g., 3000, 5173, 8000).
> 
> - Run 'hostbind context --json' for service URLs.
> - Start services with 'hostbind run --name <service>'.
> - Stop services with 'hostbind stop <service>'.
> - If something is unreachable, DO NOT guess the port. Ask HostBind.
