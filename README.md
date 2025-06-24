1. Descargar dependencias
```bash
go mod tidy
```
2. Poner la api key de gemini
3. Poner URL de supabase
4. Poner api key de supabase
5. Correr server
```bash
go run cmd/chronocode/main.go
```
6. Correr testeo de performance (tiempos de database y del agente hardcodeados)
```bash
go run performance/main.go
```
