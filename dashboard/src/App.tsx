import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import Dashboard from "./pages/Dashboard";
import ReconciliationRuns from "./pages/ReconciliationRuns";
import RunDetail from "./pages/RunDetail";
import Exceptions from "./pages/Exceptions";
import FeeAnalysis from "./pages/FeeAnalysis";
import AcquirerContracts from "./pages/AcquirerContracts";
import AgingDashboard from "./pages/AgingDashboard";

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/runs" element={<ReconciliationRuns />} />
          <Route path="/runs/:runId" element={<RunDetail />} />
          <Route path="/exceptions" element={<Exceptions />} />
          <Route path="/fees" element={<FeeAnalysis />} />
          <Route path="/contracts" element={<AcquirerContracts />} />
          <Route path="/aging" element={<AgingDashboard />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
};

export default App;
