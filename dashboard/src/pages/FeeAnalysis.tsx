import React, { useState, useEffect } from "react";
import { getFeeAnalysis } from "../api/reconciliation";
import MoneyDisplay from "../components/MoneyDisplay";

const FeeAnalysis: React.FC = () => {
  const [period, setPeriod] = useState("2024-03");
  const [analysis, setAnalysis] = useState<Record<string, unknown>>({});
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchAnalysis = async () => {
      setLoading(true);
      try {
        const data = await getFeeAnalysis(period);
        setAnalysis(data);
      } catch (error) {
        console.error("Failed to fetch fee analysis:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchAnalysis();
  }, [period]);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Fee Analysis</h1>

      {/* Period Selector */}
      <div className="bg-white rounded-lg shadow p-6">
        <label className="block">
          <span className="text-sm font-medium text-gray-700">Period (YYYY-MM)</span>
          <input
            type="month"
            value={period}
            onChange={(e) => setPeriod(e.target.value)}
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg text-gray-900"
          />
        </label>
      </div>

      {loading ? (
        <div className="text-center text-gray-500">Loading analysis...</div>
      ) : (
        <div className="space-y-6">
          {/* Summary Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Summary</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <p className="text-gray-500 text-sm">Contracted Fees</p>
                <p className="text-2xl font-bold text-gray-900 mt-1">
                  <MoneyDisplay cents={0} />
                </p>
              </div>
              <div>
                <p className="text-gray-500 text-sm">Actual Fees</p>
                <p className="text-2xl font-bold text-gray-900 mt-1">
                  <MoneyDisplay cents={0} />
                </p>
              </div>
              <div>
                <p className="text-gray-500 text-sm">Variance</p>
                <p className="text-2xl font-bold text-red-600 mt-1">
                  <MoneyDisplay cents={0} />
                </p>
              </div>
            </div>
          </div>

          {/* Detailed Breakdown */}
          <div className="bg-white rounded-lg shadow overflow-hidden">
            <div className="p-6 border-b border-gray-200">
              <h2 className="text-xl font-bold text-gray-900">By Acquirer</h2>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Acquirer</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Bandeira</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Contracted</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actual</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Variance</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  <tr className="hover:bg-gray-50">
                    <td className="px-6 py-4 text-sm font-medium text-gray-900">Sample Acquirer 1</td>
                    <td className="px-6 py-4 text-sm text-gray-600">Visa</td>
                    <td className="px-6 py-4 text-sm font-mono"><MoneyDisplay cents={0} /></td>
                    <td className="px-6 py-4 text-sm font-mono"><MoneyDisplay cents={0} /></td>
                    <td className="px-6 py-4 text-sm font-mono text-red-600"><MoneyDisplay cents={0} /></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default FeeAnalysis;
