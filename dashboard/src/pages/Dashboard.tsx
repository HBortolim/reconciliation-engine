import React, { useState, useEffect } from "react";
import { getRuns } from "../api/reconciliation";
import { ReconciliationRun } from "../types";
import MoneyDisplay from "../components/MoneyDisplay";
import StatusBadge from "../components/StatusBadge";

const Dashboard: React.FC = () => {
  const [runs, setRuns] = useState<ReconciliationRun[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRuns = async () => {
      try {
        const data = await getRuns();
        setRuns(data.slice(0, 5)); // Get last 5
      } catch (error) {
        console.error("Failed to fetch runs:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchRuns();
  }, []);

  if (loading) {
    return <div className="text-center text-gray-500">Loading...</div>;
  }

  const totalProcessed = runs.reduce((sum, r) => sum + r.totalProcessed, 0);
  const totalMatched = runs.reduce((sum, r) => sum + r.matchedCount, 0);
  const totalExceptions = runs.reduce((sum, r) => sum + r.exceptionCount, 0);
  const totalAmount = runs.reduce((sum, r) => sum + r.statistics.totalAmount, 0);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Total Processed</h3>
          <p className="text-3xl font-bold text-gray-900 mt-2">{totalProcessed}</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Matched</h3>
          <p className="text-3xl font-bold text-green-600 mt-2">{totalMatched}</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Exceptions</h3>
          <p className="text-3xl font-bold text-red-600 mt-2">{totalExceptions}</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Total Amount</h3>
          <p className="text-lg font-mono font-bold text-gray-900 mt-2">
            <MoneyDisplay centavos={totalAmount} />
          </p>
        </div>
      </div>

      {/* Recent Runs */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Recent Runs</h2>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Run ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Processed</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Matched</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Exceptions</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {runs.map((run) => (
                <tr key={run.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{run.id}</td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={run.status} variant="status" />
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{run.totalProcessed}</td>
                  <td className="px-6 py-4 text-sm text-green-600 font-medium">{run.matchedCount}</td>
                  <td className="px-6 py-4 text-sm text-red-600 font-medium">{run.exceptionCount}</td>
                  <td className="px-6 py-4 text-sm text-gray-500">{new Date(run.startTime).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
