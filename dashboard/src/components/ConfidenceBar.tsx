import React from "react";

interface ConfidenceBarProps {
  confidence: number; // 0-100
  showLabel?: boolean;
}

const getColorClass = (confidence: number): string => {
  if (confidence >= 90) return "bg-green-500";
  if (confidence >= 75) return "bg-green-400";
  if (confidence >= 60) return "bg-yellow-400";
  if (confidence >= 40) return "bg-orange-400";
  return "bg-red-500";
};

const ConfidenceBar: React.FC<ConfidenceBarProps> = ({
  confidence,
  showLabel = true
}) => {
  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 w-32 h-2 bg-gray-200 rounded-full overflow-hidden">
        <div
          className={`h-full transition-all ${getColorClass(confidence)}`}
          style={{ width: `${Math.min(confidence, 100)}%` }}
        />
      </div>
      {showLabel && (
        <span className="text-sm font-medium text-gray-700">
          {confidence.toFixed(1)}%
        </span>
      )}
    </div>
  );
};

export default ConfidenceBar;
