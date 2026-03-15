import React from "react";

interface MoneyDisplayProps {
  cents: number;
  symbol?: string;
  showCents?: boolean;
}

const MoneyDisplay: React.FC<MoneyDisplayProps> = ({
  cents,
  symbol = "R$",
  showCents = true
}) => {
  const reais = cents / 100;
  const formatted = reais.toLocaleString("pt-BR", {
    style: "currency",
    currency: "BRL",
    minimumFractionDigits: showCents ? 2 : 0,
    maximumFractionDigits: 2
  });

  return <span className="font-mono text-gray-900">{formatted}</span>;
};

export default MoneyDisplay;
