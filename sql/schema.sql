CREATE TABLE usuarios (
    id INT IDENTITY(1,1) PRIMARY KEY,
    nome NVARCHAR(255) NOT NULL,
    email NVARCHAR(255) NOT NULL UNIQUE,
    whatsapp NVARCHAR(50),
    password_hash NVARCHAR(255) NOT NULL,
    tipo_conta NVARCHAR(50) NOT NULL,
    nome_propriedade NVARCHAR(255),
    endereco_rural NVARCHAR(255),

    CONSTRAINT CK_usuarios_tipo CHECK (tipo_conta IN ('P', 'C'))
);

CREATE TABLE safras (
    id_safra INT IDENTITY(1,1) PRIMARY KEY,
    id_produtor BIGINT NOT NULL,
    nome NVARCHAR(255) NOT NULL,
    quantidade_total DECIMAL(10,2) NOT NULL DEFAULT 0,
    unidade NVARCHAR(20) NOT NULL DEFAULT 'kg',
    data_previsao DATE,
    status NVARCHAR(20) NOT NULL DEFAULT 'planejamento',
    created_at DATETIMEOFFSET NOT NULL DEFAULT SYSDATETIMEOFFSET(),

    CONSTRAINT CK_safras_status CHECK (status IN ('planejamento', 'crescimento', 'colheita', 'finalizada')),
    CONSTRAINT FK_safras_usuarios FOREIGN KEY (id_produtor) REFERENCES usuarios(id) ON DELETE CASCADE
);

CREATE TABLE produtos (
    id_produto INT IDENTITY(1,1) PRIMARY KEY,
    id_produtor BIGINT NOT NULL,
    id_safra INT,
    nome NVARCHAR(255) NOT NULL,
    descricao NVARCHAR(MAX),
    preco DECIMAL(10,2) NOT NULL DEFAULT 0,
    estoque DECIMAL(10,2) NOT NULL DEFAULT 0,
    unidade NVARCHAR(20) NOT NULL DEFAULT 'kg',
    imagem_url NVARCHAR(MAX),
    ativo BIT NOT NULL DEFAULT 1,
    created_at DATETIMEOFFSET NOT NULL DEFAULT SYSDATETIMEOFFSET(),

    CONSTRAINT FK_produtos_usuarios FOREIGN KEY (id_produtor) REFERENCES usuarios(id),
    CONSTRAINT FK_produtos_safra FOREIGN KEY (id_safra) REFERENCES safras(id_safra) ON DELETE SET NULL
);

CREATE TABLE pedidos (
    id_pedido INT IDENTITY(1,1) PRIMARY KEY,
    id_cliente BIGINT NOT NULL,
    id_produtor BIGINT NOT NULL,
    valor_total DECIMAL(10,2) NOT NULL DEFAULT 0,
    status NVARCHAR(20) NOT NULL DEFAULT 'pendente',
    observacoes NVARCHAR(MAX),
    created_at DATETIMEOFFSET NOT NULL DEFAULT SYSDATETIMEOFFSET(),

    CONSTRAINT CK_pedidos_status CHECK (status IN ('pendente', 'confirmado', 'em_preparacao', 'em_entrega', 'entregue', 'cancelado')),
    CONSTRAINT FK_pedidos_cliente FOREIGN KEY (id_cliente) REFERENCES usuarios(id),
    CONSTRAINT FK_pedidos_produtor FOREIGN KEY (id_produtor) REFERENCES usuarios(id)
);

CREATE TABLE itens_pedido (
    id_item_pedido INT IDENTITY(1,1) PRIMARY KEY,
    id_pedido INT NOT NULL,
    id_produto INT NOT NULL,
    quantidade DECIMAL(10,2) NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,

    CONSTRAINT FK_itens_pedido_pedido FOREIGN KEY (id_pedido) REFERENCES pedidos(id_pedido) ON DELETE CASCADE,
    CONSTRAINT FK_itens_pedido_produto FOREIGN KEY (id_produto) REFERENCES produtos(id_produto)
);
