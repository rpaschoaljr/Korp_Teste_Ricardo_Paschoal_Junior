-- LIMPEZA TOTAL DO BANCO
TRUNCATE TABLE itens_fatura, faturas, itens, clientes CASCADE;

-- REINICIAR SEQUENCIAS
ALTER SEQUENCE clientes_id_seq RESTART WITH 1;
ALTER SEQUENCE itens_id_seq RESTART WITH 1;
ALTER SEQUENCE faturas_id_seq RESTART WITH 1;
ALTER SEQUENCE itens_fatura_id_seq RESTART WITH 1;

-- 1. POPULAR CLIENTES (10 Clientes)
INSERT INTO clientes (nome, telefone, endereco, cpf, cnpj) VALUES
('JOAO SILVA', '11999998888', 'RUA DAS FLORES, 123', '123.456.789-01', NULL),
('MARIA OLIVEIRA', '21988887777', 'AVENIDA CENTRAL, 500', '234.567.890-12', NULL),
('KORP SOLUCOES', '4133332222', 'RUA TECNOLOGICA, 1000', NULL, '12.345.678/0001-90'),
('TECH SOLUTIONS', '1144445555', 'AV PAULISTA, 1500', NULL, '98.765.432/0001-10'),
('CARLOS SOUZA', '31977776666', 'RUA MINAS, 45', '345.678.901-23', NULL),
('ANA COSTA', '51966665555', 'AV SUL, 90', '456.789.012-34', NULL),
('REDE VAREJO', '1122223333', 'RUA DO COMERCIO, 10', NULL, '11.222.333/0001-44'),
('DISTRIBUIDORA NORTE', '8134445555', 'AV PERIMETRAL, 200', NULL, '55.666.777/0001-88'),
('PAULO SANTOS', '11955554444', 'RUA JARDIM, 8', '567.890.123-45', NULL),
('BEATRIZ LIMA', '19944443333', 'RUA CAMPINAS, 300', '678.901.234-56', NULL);

-- 2. POPULAR ITENS (100 Itens)
DO $$
BEGIN
    FOR i IN 1..100 LOOP
        INSERT INTO itens (codigo, descricao, saldo, preco_base)
        VALUES (
            'PROD-' || LPAD(i::text, 3, '0'),
            'PRODUTO TESTE ' || LPAD(i::text, 3, '0'),
            floor(random() * 200 + 10)::int, -- Saldo entre 10 e 210
            (random() * 500 + 5.0)::decimal(10,2) -- Preço entre 5 e 505
        );
    END LOOP;
END $$;

-- 3. POPULAR FATURAS (100 Notas)
-- Espalhadas nos últimos 30 dias com horas variadas
DO $$
DECLARE
    f_id int;
    c_id int;
    total_val decimal(10,2);
    data_fat timestamp;
    status_fat varchar(20);
BEGIN
    FOR i IN 1..100 LOOP
        c_id := floor(random() * 10 + 1)::int; -- Cliente aleatório entre 1 e 10
        data_fat := NOW() - (random() * interval '30 days'); -- Data aleatória nos ultimos 30 dias
        
        -- 80% das notas serão FECHADAS para termos dados para a IA analisar
        IF random() > 0.2 THEN
            status_fat := 'FECHADA';
        ELSE
            status_fat := 'ABERTA';
        END IF;

        INSERT INTO faturas (cliente_id, status, valor_total, data_criacao)
        VALUES (c_id, status_fat, 0, data_fat)
        RETURNING id INTO f_id;

        -- Adicionar 1 a 5 itens por fatura
        total_val := 0;
        FOR j IN 1..(floor(random() * 5 + 1)::int) LOOP
            DECLARE
                item_id_rand int;
                qtd_rand int;
                prec_rand decimal(10,2);
                sub_total decimal(10,2);
            BEGIN
                item_id_rand := floor(random() * 100 + 1)::int; -- Produto aleatório
                qtd_rand := floor(random() * 10 + 1)::int; -- Quantidade 1-10
                
                SELECT preco_base INTO prec_rand FROM itens WHERE id = item_id_rand;
                sub_total := qtd_rand * prec_rand;
                total_val := total_val + sub_total;

                INSERT INTO itens_fatura (fatura_id, item_id, quantidade, preco_unitario, subtotal)
                VALUES (f_id, item_id_rand, qtd_rand, prec_rand, sub_total);
            END;
        END LOOP;

        -- Atualizar o total da fatura
        UPDATE faturas SET valor_total = total_val WHERE id = f_id;

    END LOOP;
END $$;

-- 4. FINALIZAR: Ajustar o saldo de alguns produtos para simular alerta de estoque na IA
UPDATE itens SET saldo = 2 WHERE codigo = 'PROD-001';
UPDATE itens SET saldo = 1 WHERE codigo = 'PROD-005';
UPDATE itens SET saldo = 0 WHERE codigo = 'PROD-010';
