import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress as ProgressBar } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { BookOpen, Clock, Target, TrendingUp, Calendar } from "lucide-react";
import { useNavigate } from "react-router-dom";

interface Essay {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  wordCount: number;
}

interface StudyStats {
  essayId: string;
  title: string;
  studyCount: number;
  reciteCount: number;
  dictationCount: number;
  lastStudied: string;
  completionRate: number;
}

const Progress = () => {
  const navigate = useNavigate();
  const [essays, setEssays] = useState<Essay[]>([]);
  const [studyStats, setStudyStats] = useState<StudyStats[]>([]);
  const [totalStats, setTotalStats] = useState({
    totalEssays: 0,
    totalWords: 0,
    totalStudyTime: 0,
    avgCompletionRate: 0
  });

  useEffect(() => {
    const essaysData = JSON.parse(localStorage.getItem('essays') || '[]');
    setEssays(essaysData);

    // Generate mock study statistics
    const mockStats: StudyStats[] = essaysData.map((essay: Essay) => ({
      essayId: essay.id,
      title: essay.title,
      studyCount: Math.floor(Math.random() * 10) + 1,
      reciteCount: Math.floor(Math.random() * 8) + 1,
      dictationCount: Math.floor(Math.random() * 5) + 1,
      lastStudied: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
      completionRate: Math.floor(Math.random() * 40) + 60
    }));

    setStudyStats(mockStats);

    // Calculate total statistics
    const totalWords = essaysData.reduce((sum: number, essay: Essay) => sum + essay.wordCount, 0);
    const avgCompletion = mockStats.length > 0 
      ? mockStats.reduce((sum, stat) => sum + stat.completionRate, 0) / mockStats.length 
      : 0;

    setTotalStats({
      totalEssays: essaysData.length,
      totalWords,
      totalStudyTime: mockStats.reduce((sum, stat) => sum + stat.studyCount * 15, 0), // 15 min per study
      avgCompletionRate: avgCompletion
    });
  }, []);

  const getStudyLevel = (studyCount: number) => {
    if (studyCount >= 8) return { level: "Expert", color: "bg-primary" };
    if (studyCount >= 5) return { level: "Advanced", color: "bg-accent" };
    if (studyCount >= 3) return { level: "Intermediate", color: "bg-secondary" };
    return { level: "Beginner", color: "bg-muted" };
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-foreground mb-2">Learning Progress</h1>
        <p className="text-muted-foreground">Track your English essay learning journey</p>
      </div>

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center space-x-2">
              <BookOpen className="h-4 w-4 text-primary" />
              <div>
                <p className="text-2xl font-bold">{totalStats.totalEssays}</p>
                <p className="text-xs text-muted-foreground">Total Essays</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center space-x-2">
              <Target className="h-4 w-4 text-accent" />
              <div>
                <p className="text-2xl font-bold">{totalStats.totalWords.toLocaleString()}</p>
                <p className="text-xs text-muted-foreground">Words Studied</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center space-x-2">
              <Clock className="h-4 w-4 text-primary" />
              <div>
                <p className="text-2xl font-bold">{Math.round(totalStats.totalStudyTime / 60)}h</p>
                <p className="text-xs text-muted-foreground">Study Time</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center space-x-2">
              <TrendingUp className="h-4 w-4 text-accent" />
              <div>
                <p className="text-2xl font-bold">{Math.round(totalStats.avgCompletionRate)}%</p>
                <p className="text-xs text-muted-foreground">Avg Completion</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Individual Essay Progress */}
      <Card>
        <CardHeader>
          <CardTitle>Essay Progress Details</CardTitle>
        </CardHeader>
        <CardContent>
          {studyStats.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground mb-4">No study data available yet</p>
              <Button onClick={() => navigate('/import')}>
                Import Your First Essay
              </Button>
            </div>
          ) : (
            <div className="space-y-4">
              {studyStats.map((stat) => {
                const studyLevel = getStudyLevel(stat.studyCount);
                return (
                  <div key={stat.essayId} className="border rounded-lg p-4">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex-1">
                        <h3 className="font-semibold text-foreground">{stat.title}</h3>
                        <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                          <span className="flex items-center gap-1">
                            <Calendar className="h-3 w-3" />
                            Last studied: {formatDate(stat.lastStudied)}
                          </span>
                          <Badge variant="secondary" className={studyLevel.color}>
                            {studyLevel.level}
                          </Badge>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-medium">{stat.completionRate}% complete</p>
                        <ProgressBar value={stat.completionRate} className="w-24 mt-1" />
                      </div>
                    </div>

                    <div className="grid grid-cols-3 gap-4 text-sm">
                      <div className="text-center p-2 bg-muted/50 rounded">
                        <p className="font-semibold text-primary">{stat.studyCount}</p>
                        <p className="text-muted-foreground">Study Sessions</p>
                      </div>
                      <div className="text-center p-2 bg-muted/50 rounded">
                        <p className="font-semibold text-accent">{stat.reciteCount}</p>
                        <p className="text-muted-foreground">Recite Practice</p>
                      </div>
                      <div className="text-center p-2 bg-muted/50 rounded">
                        <p className="font-semibold text-primary">{stat.dictationCount}</p>
                        <p className="text-muted-foreground">Dictation Tests</p>
                      </div>
                    </div>

                    <div className="flex gap-2 mt-3">
                      <Button 
                        size="sm" 
                        onClick={() => navigate(`/essay/${stat.essayId}`)}
                      >
                        View Essay
                      </Button>
                      <Button 
                        size="sm" 
                        variant="outline"
                        onClick={() => navigate(`/essay/${stat.essayId}/recite`)}
                      >
                        Practice
                      </Button>
                      <Button 
                        size="sm" 
                        variant="outline"
                        onClick={() => navigate(`/essay/${stat.essayId}/dictation`)}
                      >
                        Dictation
                      </Button>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default Progress;