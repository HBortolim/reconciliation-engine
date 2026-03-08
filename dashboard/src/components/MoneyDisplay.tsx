import React from "react";

interface MoneyDisplayProps {
  centavos: number;
  symbol?: string;
  showCents?: boolean;
}

const MoneyDisplay: React.FC<MoneyDisplayProps> = ({
  centavos,
  symbol = "R$",
  showCents = true
}) => {
  const reais = centavos / 100;
  const formatted = reais.toLocaleString("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: showCents ? 2 : 0,
    maximumFractionDigits: 2
  });

  return <span className="font-mono text-gray-900">{formatted}</span>;
};

export default MoneyDisplay;
