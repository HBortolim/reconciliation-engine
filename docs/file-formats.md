# Supported File Formats

## Overview

The Reconciliation Engine ingests payment data from multiple Brazilian financial sources. Each source produces files in different formats with different identification keys, settlement windows, and fee structures.

---

## Bank Statements

### OFX (Open Financial Exchange)

- **Versions**: OFX 1.x (SGML — not XML), OFX 2.x (XML)
- **Source**: Bank account statement exports
- **Key fields**: STMTTRN entries (TRNTYPE, DTPOSTED, TRNAMT, FITID, MEMO)
- **Gotcha**: OFX 1.x uses SGML syntax with unclosed tags — standard XML parsers will fail
- **Banks**: All major Brazilian banks export OFX

### CNAB 240

- **Standard**: Febraban (Federação Brasileira de Bancos)
- **Format**: Fixed-width positional, 240 characters per line
- **Structure**: Header de arquivo → Header de lote → Detalhe(s) → Trailer de lote → Trailer de arquivo
- **Key fields**: Tipo de registro (pos 8), Segmento (pos 14 for detalhe), Valor (variable by bank)
- **Gotcha**: Bank-specific field variations within the "standard" — field positions for valor, data, and identificação change per bank
- **Supported banks**: Itaú, Bradesco, Banco do Brasil, Santander, Caixa, BTG Pactual, Inter, Nubank, Sicredi, Sicoob

### CNAB 400

- **Standard**: Febraban (legacy)
- **Format**: Fixed-width positional, 400 characters per line
- **Structure**: Header → Detalhe(s) → Trailer
- **Usage**: Still common for boleto remessa/retorno files
- **Gotcha**: Even more bank variation than CNAB 240; some banks use non-standard record types

---

## Pix Settlement

- **Format**: JSON or CSV exports from PSP (Payment Service Provider) dashboards
- **Key fields**: EndToEndId (E2EID — 32 chars, starts with 'E'), valor, data liquidação, chave Pix, tipo (instant/agendado/cobrança/devolução)
- **E2EID format**: `E{ISPB_8chars}{YYYYMMDD}{SEQUENCE_11chars}`
- **Special cases**: Pix Devolução (partial/full), Pix Cobrança (with txid), Pix Troco/Saque

---

## Card Acquirers

### CIELO

- **Formats**: EEFI (Extrato Eletrônico de Fluxo de Caixa), EEVC (Extrato Eletrônico de Vendas com Plano de Pagamento)
- **Type**: Fixed-width positional
- **Key fields**: NSU, Valor bruto, Valor líquido, Taxa MDR, Bandeira, Parcela, Previsão de pagamento
- **Settlement**: D+1 débito, D+30 crédito (or D+2 with antecipação)

### Rede

- **Format**: CSV
- **Key fields**: NSU, Data venda, Data pagamento, Valor bruto, Valor líquido, Taxa, Bandeira, Modalidade
- **Encoding**: ISO-8859-1 (Latin-1) — common gotcha

### Stone

- **Format**: JSON (API response)
- **Key fields**: NSU, stone_code, amount, net_amount, fee, card_brand, installments, payment_date
- **API**: RESTful with API key authentication

### PagSeguro

- **Format**: CSV
- **Key fields**: Código da transação, Data, Valor bruto, Valor líquido, Taxa PagSeguro, Tipo de pagamento
- **Encoding**: UTF-8

### Getnet

- **Format**: Fixed-width positional
- **Key fields**: NSU, Valor bruto, Valor líquido, Taxa, Bandeira, Tipo operação
- **Settlement**: Variable per contract

### SafraPay

- **Format**: CSV or fixed-width (varies by report type)
- **Key fields**: NSU, Valor, Taxa, Bandeira, Data prevista
- **Notes**: Newer acquirer, format may evolve

---

## Common Challenges

1. **Encoding**: Mix of UTF-8, ISO-8859-1, and Windows-1252 across different sources
2. **Date formats**: DD/MM/YYYY (Brazilian standard), YYYYMMDD (CNAB), ISO 8601 (APIs)
3. **Decimal separators**: Comma (,) in Brazilian locale vs. dot (.) in APIs
4. **Amount representation**: Some files use cents (integer), others use reais (decimal)
5. **Line endings**: Mix of CRLF (Windows) and LF (Unix) in bank files
6. **Empty/malformed records**: Banks sometimes produce files with blank lines or truncated records
7. **Character encoding in names**: Razão social / nome fantasia with accented characters (ç, ã, é, etc.)
