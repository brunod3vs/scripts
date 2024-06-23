package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Função para verificar a presença de strings interessantes
func checkInterestingContent(content string) {
	fmt.Println("Verificando conteúdo interessante...")

	// Comentários que podem conter informações sensíveis
	commentRegex := regexp.MustCompile(`<!--([\s\S]*?)-->`)
	comments := commentRegex.FindAllString(content, -1)
	if len(comments) > 0 {
		for _, comment := range comments {
			// Filtrar comentários vazios ou irrelevantes
			cleanComment := strings.TrimSpace(comment)
			if cleanComment != "<!-- -->" && cleanComment != "<!---->" && cleanComment != "<!--[-->" && cleanComment != "<!--[if IE]-->" && len(cleanComment) > 10 {
				fmt.Println("Comentário encontrado:", comment)
			}
		}
	} else {
		fmt.Println("Nenhum comentário relevante encontrado.")
	}

	// Scripts externos
	scriptRegex := regexp.MustCompile(`<script[^>]*src=["']([^"']*)["'][^>]*>`)
	scripts := scriptRegex.FindAllStringSubmatch(content, -1)
	if len(scripts) > 0 {
		for _, script := range scripts {
			fmt.Println("Script externo encontrado:", script[1])
		}
	} else {
		fmt.Println("Nenhum script externo encontrado.")
	}

	// Metadados
	metaRegex := regexp.MustCompile(`<meta[^>]+content=["']([^"']*)["'][^>]*>`)
	metas := metaRegex.FindAllStringSubmatch(content, -1)
	if len(metas) > 0 {
		for _, meta := range metas {
			fmt.Println("Metadado encontrado:", meta[1])
		}
	} else {
		fmt.Println("Nenhum metadado encontrado.")
	}
}

// Função para extrair subdomínios
func extractSubdomains(content string, baseDomain string) []string {
	fmt.Println("Extraindo subdomínios...")

	subdomainRegex := regexp.MustCompile(`https?://([a-zA-Z0-9.-]+)\.` + regexp.QuoteMeta(baseDomain))
	matches := subdomainRegex.FindAllStringSubmatch(content, -1)
	subdomains := make(map[string]bool)
	for _, match := range matches {
		subdomains[match[1]] = true
	}

	var uniqueSubdomains []string
	for subdomain := range subdomains {
		uniqueSubdomains = append(uniqueSubdomains, subdomain)
	}

	if len(uniqueSubdomains) > 0 {
		fmt.Println("Subdomínios encontrados:")
		for _, subdomain := range uniqueSubdomains {
			fmt.Println(subdomain)
		}
	} else {
		fmt.Println("Nenhum subdomínio encontrado.")
	}

	return uniqueSubdomains
}

// Função para buscar a URL e analisar o conteúdo
func analyzeURL(targetURL string, baseDomain string) {
	fmt.Println("Fazendo requisição para:", targetURL)

	// Configure o cliente HTTP com um tempo de espera maior
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Crie uma requisição HTTP
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	// Adicione cabeçalhos para imitar um navegador real
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Referer", "https://www.google.com")

	// Envie a requisição
	resp, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Erro: A conexão expirou.")
		} else {
			fmt.Println("Erro ao fazer requisição:", err)
		}
		return
	}
	defer resp.Body.Close()

	fmt.Println("Resposta recebida, status code:", resp.StatusCode)

	// Verifica se a resposta é 403
	if resp.StatusCode == http.StatusForbidden {
		fmt.Println("Acesso proibido (403). Não foi possível analisar esta URL.")
		return
	}

	// Leia o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return
	}

	fmt.Println("Resposta lida, tamanho:", len(body))

	// Converta o corpo para string
	content := string(body)

	// Verifique o conteúdo em busca de informações interessantes
	checkInterestingContent(content)

	// Extrair e analisar subdomínios
	subdomains := extractSubdomains(content, baseDomain)
	for _, subdomain := range subdomains {
		// Adicionar um intervalo aleatório entre as requisições
		time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)
		fmt.Println("Verificando subdomínio:", subdomain)
		analyzeURL("https://"+subdomain, baseDomain)
	}
}

func main() {
	// Perguntar ao usuário qual URL ele quer verificar
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite a URL que você quer verificar: ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL)

	// Garantir que a URL inclui um esquema
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	// Extrair domínio base
	u, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("URL inválida:", err)
		return
	}
	baseDomain := u.Hostname()

	fmt.Println("Domínio base extraído:", baseDomain)

	// Analisar a URL fornecida
	analyzeURL(targetURL, baseDomain)
}
