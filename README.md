# Magalu Cloud Back-end Desafio

## Parte 1:

Como os princípios SOLID e de design funcional podem ser aplicados para maximizar a capacidade das equipes de integrar, estender e evoluir arquiteturas complexas de software?

R: Os princípios SOLID e o design funcional podem oferecer uma arquitetura de software modular, flexível e testável. Basicamente, o SOLID nos auxilia na divisão de responsabilidades das nossas implementações. Por exemplo, o SRP (Single Responsibility Principle) afirma que uma função ou classe deve ter apenas uma responsabilidade. Somado à DIP (Dependency Inversion Principle), isso possibilita alterar nossa implementação utilizando interfaces ou camadas de abstração (Interface Segregation Principle). Dessa forma, injetamos nossa implementação no construtor da classe.

Utilizando o ISP (Interface Segregation Principle), criamos uma camada que separa nossa implementação da abstração, facilitando alterações em frameworks ou regras de negócio. Essas práticas também colaboram para a aplicação do OCP (Open/Closed Principle), que dita que uma implementação deve estar aberta para extensão, porém fechada para alteração. Com essas estratégias, o SOLID contribui para a segregação de responsabilidades, facilitando os testes e a leitura do código, além de tornar claras as camadas de negócio e de infraestrutura, por exemplo.

O design funcional, por sua vez, nos ajuda a dividir grandes implementações em pequenas funções, o que melhora a legibilidade do código e torna mais fácil identificar onde alterar ou implementar novas funcionalidades. A utilização da composição de funções permite criar uma classe ou função responsável por gerenciar pequenas funções. Dessa maneira, além de aplicar o OCP, ICP e o SRP, contribuímos para a manutenibilidade e o registro (logging) do nosso código, sendo muito mais simples debugar uma função pequena de 50 linhas com seu contexto delimitado do que navegar por classes extensas com centenas de linhas de código.

## Parte 2:

Nesta parte do desafio, implementei a lógica de Ingestor, para isso levei alguns pontos na decisão da arquitetura:

- Fault Tolarance: Como podemos garantir caso alguns desses micro serviços fique fora do ar, seja possível reprocessar todas mensagens a fim que nenhum "pulso" seja perdido da contabilidade.

- Escala: Como podemos escalar verticalmente a infraestrutura para que ela suporte a demanda crescente de requisições e processamento.

![Logo](assets/software_architeture.png)

### Overview

### Ingestor Api:
Decisões: 
gRPC: A serialização do gRPC é via Protocol Buffers, que por ser em binários tende a ser de 7 a 10 vezes mais rápido que HTTP/HTTPS, melhorando a performance de I/O de nossa API, sem necessidade de escalar até que atinja um número bem significativos de requisições. Sendo que o contrato para nossa API tende a ser centralizado em proto files, não teremos problemas de clientes mandando interfaces erradas ou que sejam dificies de serem mantidas conforme vamos incrementando ela.

Kafka: Após receber nosso "pulso", vamos enviar diretamente para uma fila do Kafka, para que esses dados sejam processados de forma assíncrona, visto que o client não necessita dessas informações tratadas, não neste momento.

### Ingestor Consumer:
Decisões:
Kafka: Teremos uma Job que pode ser controlada de quanto em quanto tempo irá fazer a leitura no tópico `resource-consumption`, durante o certo período vamos ler todas mensagens e agregalas. O motivo de usar um consumer que possui uma quantidade especifica que ficara online, será pelo motivo de pulling do kafka, vamos processando em batches os dados e disponibilizando sob demanda para o `Processador e Armazenador`. Sempre que um consumer fica offline, o kafka salva as mensagens como "Messages Behind", e sempre que o consumer subir novamente, vamos processar esses dados.

Go Channel: Não processamos diretamente do consumer os dados, ao invés disso, pegamos a mensagem e redirecionamos para um channel, para ser processado em uma goroutine a parte, isso nos ajuda a desacoplar nossa aplicação do kafka, sendo assim podemos ler as mensagens de RabbitMQ, Redis Queue, por exemplo, desde que a mensagem seja enviada para dentro dessa goroutine.

### Processador e Armazenador:
Decisões:
Nessa parte podemos seguir de algumas formas, dependendo de nossas necessidades. 
Caso 01: Caso nós precisemos que seja falha a tolerância, como por exemplo a criação de uma DeadLetter queue, caso nossa aplicação tenha alguma falha, podemos redirecionar para essa outra fila e fazer o processamento retroativo posteriormente.
![Logo](assets/ConsumerContratos/architeture_2.png)

Caso 02: Caso a tolerância a falha não seja necessária no lado do nosso consumer e sim o retry por parte do client, podemos usar uma comunicação direta por gRPC, facilitando nosso fluxo de comunicação. O único problema nesse approach é caso surja necessidade de reprocessamento de mensagens anteriores, vamos ter necessidade de força o re-envio a aprtir do Ingestor. Ainda sim, manteria todos os dados salvos como backup dentro de um banco de dados feito para alto I/O como ScyllaDB ou CockroachDB. Eles poderiam ser utilizados para extrassão de analytics posteriormente caso encessário
![Logo](assets/ConsumerContratos/architeture_1.png)

### Contratos e Catálogo:
Decisões:
Catálogo deverá consumir os dados agregados, e precificar corretamente, o input poderia ser tanto com um consumer em um broker (como kafka), ou um entry point como HTTP ou gRPC.

### Consulta:
Decisões:
Visto que essa API tem o propósito de servir a um front-end, seria mais interessante implementa-lá em Python, com o uso do FastAPI e a geração automática de documentação no OpenAPI ou Swagger, iria facilitar a consulta dos nossos clients, nesse contexto, gRPC não seria muito útil, visto que ainda é uma tecnologia de uso emergente para web browser, e podem ter problemas de compatibilidade.