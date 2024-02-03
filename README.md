# Execução e testes via browser
## Para subir o sistema: 
Fazer o build da imagem:

    docker compose build --no-cache

Subir a imagem co o comando:

    docker compose up -d   

O servidor estará exposto na porta 8080

## Para realizar as chamadas acesse o seguinte endereço:
    http://localhost:8080/

No arquivo 

    request.http 

existem alguns exemplos de acessos com diferentes configurações de tokens e IPs.<br>
Executando cada uma delas repedidas vezes irá ativar as validações do Rate Limiter.

<br>

# Definição dos tempos e limites de acesso
As restrições de acesso são aplicadas por padrão aos IPs de quem faz a requisição. Opcionalmente a restrição pode ser aplciada a tokens, caso esses estejam presentes no request conforme exemplo abaixo (no formato para um arquivo request.http):

    GET http://localhost:8080/ HTTP/1.1
    Host: localhost:8000
    API_KEY: token2

3 tipos de restrições são aplicadas
- Número máximo de acessos: Número de vezes que um IP ou Token pode realizar requisições dentro de um tempo limite. Ultrapassando este limite, o IP ou Token é bloqueado;
- Tempo limite: Tempo no qual será contabilizado o número de acessos. É definido em segundos;
- Tempo de Bloqueio: Tempo, em segundos, que um IP ou Token ficará bloqueado caso ultrapasse o número de acessos permitido no tempo limite.

As configurções de IP são sempre genéricas e defindas no arquivo .env (ver mais abaixo).
As configurações de token podem ser genérias, também defindas no arquivo **.env**, ou especificadas por token. Neste caso os limites de cada token devem ser definidos no arquivo definido na propriedade TOKEN_FILE_LIMITS (ver mais abaixo).

# Configuração das restrições no arquivo **.env**
## Configurações para IP
As 3 propriedades abaixo precisam estar definidas e serão aplicadas individualmente a cada IP de **todas** as requisições que **não** tenham um token especificado:

	IP_MAX_NUMBER_ACCESS: Número máximo de acessos
	IP_TIME_LIMIT: Tempo limite
	IP_TIME_BLOCK: Tempo de Bloqueio

Exemplo de preenchimento

	IP_MAX_NUMBER_ACCESS=5
	IP_TIME_LIMIT=5
	IP_TIME_BLOCK=10

## Configurações para Token
As 4 propriedades abaixo precisam estar definidas e serão aplicadas individualmente a cada Token de das requisições que **tenham** um token especificado, e o token não esteja presente no arquivo de tokens específicos:

	TOKEN_MAX_NUMBER_ACCESS: Número máximo de acessos
	TOKEN_TIME_LIMIT: Tempo limite
	TOKEN_TIME_BLOCK: Tempo de Bloqueio
	TOKEN_FILE_LIMITS: arquivo onde estarão definidas as retrições específicas de cada token

Exemplo de preenchimento

	TOKEN_MAX_NUMBER_ACCESS=5
	TOKEN_TIME_LIMIT=5
	TOKEN_TIME_BLOCK=10
	TOKEN_FILE_LIMITS=tokens.json

O arquivo definido na propriedade TOKEN_FILE_LIMITS deve definir os seguintea parâmetros:

    token: identificação do token
    maxNumberAccess: Número máximo de acessos
    timeLimit: Tempo limite
    timeBlock: Tempo de Bloqueio

Abaixo segue um exemplo de preenchimento do arquivo:

    [
    {
        "token": "token7",
        "maxNumberAccess": 6,
        "timeLimit": 20,
        "timeBlock": 10
    },
    {
        "token": "token8",
        "maxNumberAccess": 10,
        "timeLimit": 7,
        "timeBlock": 11
    }  
    ]

# Como é feito o controle de acessos

A cada acesso o rate limiter:

- Verifica se existe um registro de bloqueio para o IP/token. Neste caso a resposta é configurada para indicar que o número de acessos foi atingida e **o fluxo da requisição é interrompido.**

- verifica se já existe registro de acesso do IP/token. Caso não exista é criado um registro no Redis, incializado com 1 (primeiro acesso) e com o tempo de expiração definido ou no arquivo **.env**, ou no arquivo **tokens.json**. **O fluxo da requisição segue normalmente.**

- Se já existe registro no redis, valida o número de acessos. Se este valor for igual ou maior ao do número máximo de acessos permitido, um registro de bloqueio é incluído no redis com uma marcação específica indicando que aquele IP ou token foi bloqueado. Este registro é configurado como tempo de expiração definido no Tempo de Bloqueio e o Redis se encarrega de excluí-lo quando o tempo expirar. A resposta é configurada para indicar que o número de acessos foi atingida e **o fluxo da requisição é interrompido.**

- na última etapa, caso chegue nela, o número de acessos do IP/token é incrementado e **o fluxo da requisição segue normalmente.**

O redis controla o tempo de expiração dos registros, e quando este tempo é atingido o registro é automaticamente excluído, "zerando" o contador de acessos daquele período. Um novo acesso cria um novo registro e o ciclo continua.


# Implementação

## main.go
Logo no início é feita a leitura das propriedades do **.env**
Neste arquivo estão contidas:
- as configurações do servidor que irá guardar os registros de acessos. Daqui para frente chamdo de **BD**
- Configurações de restrições por IP
- Configurações de restrições por token

Em seguida um factory de BD irá criar um instância de BD de acordo com o tipo configurado (atualmente suporta apenas o Redis). Esta instâcia implementa a interface de BD o que permite alterar o tipo de BD sem que o core tenha que ser reescrito. 

É criada a instâcia da **rateLimiter** que tem as regras de negócio a serem aplicadas aos acessos.

É feita a limpeza da lista de restrições de tokens no BD e em seguida a nova lista é carregada do arquivo especificado na propriedade **TOKEN_FILE_LIMITS** e inserida no BD.

É criada um instância de roteamento e nela é injetada a função **ratelimit** que implementa o que é necessário para que atue como midleware e filtra as chamadas para que possam ser avaliadas pelo Rate Limiter.
Nesta função **ratelimit** são coletadas as características da requisição que são passadas para o método **TrataRatelimit** da instancia do **retelimiter**, que fará então a validação dos acessos antes que o request seja completado.

O servidor http é iniciado.

## internal/entity
### rateLimiter.go

Esta classe implementas as regras de negócio do limitador de acessos, e nçao depende de nenhum tipo de BD específico ou qualquer outro recurso.
Tem como propriedades um BD que deve implementar a interface de BD e os valores de limites que dem ser aplicados.

O método principal é o **TrataRatelimit** e nele é implementada todas as regrajá descritas no tópico **Como é feito o controle de acessos** acima. 
Caso o limite seja por IP são utilizados os valores padrão de restriçõe de IP passados como parâmetro na construção da instância.
Caso haja um token na requisição, é verificado se existe um regisoto de limites para este token no BD. Caso existe o valores deste registro serão utilizados. Caso não exista são utilizados os valores padrão de restriçõe de token passados como parâmetro na construção da instância.

Esta classe também implementa métodos que encapsulam as chamadas de BD, o que permite os testes unitários.

## internal/database
### interfaces.go

Define os métodos que precisam ser implementados pelas classes que representam algum tipo de banco de dados, para atender às regra de negócio aplicadas às limitações de tráfico 

### RedisRL
Implementa o BD do tipo Redis
Possui a implementação de todos os métodos definidos na interface acessando o Redis


## Testes
Todas as principais classes possuem testes unitários