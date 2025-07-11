import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Header } from "@/components/layout/Header";
import { Footer } from "@/components/layout/Footer";
import Home from "./pages/Home";
import Import from "./pages/Import";
import Library from "./pages/Library";
import EssayDetail from "./pages/EssayDetail";
import ReciteMode from "./pages/ReciteMode";
import Dictation from "./pages/Dictation";
import Progress from "./pages/Progress";
import ClozeTest from "./pages/ClozeTest";
import TagsPage from "./pages/TagsPage";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <div className="min-h-screen flex flex-col">
          <Header />
          <main className="flex-1">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/import" element={<Import />} />
              <Route path="/library" element={<Library />} />
              <Route path="/essay/:id" element={<EssayDetail />} />
              <Route path="/essay/:id/recite" element={<ReciteMode />} />
              <Route path="/essay/:id/dictation" element={<Dictation />} />
              <Route path="/essay/:id/cloze" element={<ClozeTest />} />
              <Route path="/progress" element={<Progress />} />
              <Route path="/tags" element={<TagsPage />} />
              <Route path="/tags/:tag" element={<TagsPage />} />
              {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
              <Route path="*" element={<NotFound />} />
            </Routes>
          </main>
          <Footer />
        </div>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
