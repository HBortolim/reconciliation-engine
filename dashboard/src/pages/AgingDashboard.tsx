import React, { useState, useEffect } from "react";
import MoneyDisplay from "../components/MoneyDisplay";

interface AgingBucket {
  label: string;
  range: string;
  count: number;
  amount: number;
}

const AgingDashboard: React.FC = () => {
  const [buckets, setBuckets] = useState<AgingBucket[]>([
    { label: "0-7 days", range: "0-7", count: 45, amount: 125000 },
    { label: "8-15 days", range: "8-15", count: 32, amount: 87500 },
    { label: "16-30 days", range: "16-30", count: 18, amount: 52300 },
    { label: "30+ days", range: "30+", count: 8, amount: 23100 }
  ]);

  const totalCount = buckets.reduce((sum, b) => sum + b.count, 0);
  const totalAmount = buckets.reduce((sum, b) => sum + b.amount, 0);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Aging Report</h1>

      {/* Summary */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Total Unreconciled Items</h3>
          <p className="text-4xl font-bold text-gray-900 mt-2">{totalCount}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Total Unreconciled Amount</h3>
          <p className="text-2xl font-mono font-bold text-gray-900 mt-2">
            <MoneyDisplay centavos={totalAmount} />
          </p>
        </div>
      </div>

      {/* Aging Buckets */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {buckets.map((bucket) => (
          <div key={bucket.range} className="bg-white rounded-lg shadow p-6">
            <h3 className="text-gray-500 text-sm font-medium">{bucket.label}</h3>
            <div className="mt-4 space-y-3">
              <div>
                <p className="text-2xl font-bold text-gray-900">{bucket.count}</p>
                <p className="text-xs text-gray-500">Items</p>
              </div>
              <div>
                <p className="text-lg font-mono font-bold text-gray-900">
                  <MoneyDisplay centavos={bucket.amount} />
                </p>
                <p className="text-xs text-gray-500">Amount</p>
              </div>
              <div className="pt-2 border-t border-gray-200">
                <p className="text-sm text-gray-600">
                  {((bucket.count / totalCount) * 100).toFixed(1)}% of total
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Detailed Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Unreconciled Items by Age</h2>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Age Bucket</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Count</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">% of Total</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {buckets.map((bucket) => (
                <tr key={bucket.range} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{bucket.label}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{bucket.count}</td>
                  <td className="px-6 py-4 text-sm text-gray-900 font-mono">
                    <MoneyDisplay centavos={bucket.amount} />
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {((bucket.count / totalCount) * 100).toFixed(1)}%
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <button className="text-blue-600 hover:text-blue-800 font-medium">View Items</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default AgingDashboard;
