import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { 
  Play, Pause, SkipForward, SkipBack, Eye, EyeOff, 
  ArrowLeft, RotateCcw, CheckCircle, Volume2 
} from "lucide-react";
import { useToast } from "@/hooks/use-toast";

interface Essay {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  audioUrl: string;
  imageUrl: string;
  wordCount: number;
}

export default function ReciteMode() {
  const { id } = useParams<{ id: string }>();
  const [essay, setEssay] = useState<Essay | null>(null);
  const [sentences, setSentences] = useState<string[]>([]);
  const [currentSentenceIndex, setCurrentSentenceIndex] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [showText, setShowText] = useState(false);
  const [completedSentences, setCompletedSentences] = useState<Set<number>>(new Set());
  const { toast } = useToast();

  useEffect(() => {
    if (id) {
      const savedEssays = JSON.parse(localStorage.getItem('essays') || '[]');
      const foundEssay = savedEssays.find((e: Essay) => e.id === id);
      if (foundEssay) {
        setEssay(foundEssay);
        // 将文章分割成句子
        const sentenceArray = foundEssay.content
          .split(/[.!?]+/)
          .map(s => s.trim())
          .filter(s => s.length > 0);
        setSentences(sentenceArray);
      }
    }
  }, [id]);

  const playCurrentSentence = () => {
    setIsPlaying(true);
    // 模拟播放当前句子
    toast({
      title: "播放当句子音频",
      description: `正在播放第 ${currentSentenceIndex + 1} 句`,
    });
    
    // 模拟音频播放时长
    setTimeout(() => {
      setIsPlaying(false);
    }, 2000);
  };

  const nextSentence = () => {
    if (currentSentenceIndex < sentences.length - 1) {
      setCurrentSentenceIndex(currentSentenceIndex + 1);
      setShowText(false);
      setIsPlaying(false);
    }
  };

  const previousSentence = () => {
    if (currentSentenceIndex > 0) {
      setCurrentSentenceIndex(currentSentenceIndex - 1);
      setShowText(false);
      setIsPlaying(false);
    }
  };

  const toggleTextVisibility = () => {
    setShowText(!showText);
  };

  const markCurrentAsCompleted = () => {
    const newCompleted = new Set(completedSentences);
    newCompleted.add(currentSentenceIndex);
    setCompletedSentences(newCompleted);
    
    toast({
      title: "已标记为完成",
      description: "这句话已经掌握了！",
    });
  };

  const resetProgress = () => {
    setCurrentSentenceIndex(0);
    setCompletedSentences(new Set());
    setShowText(false);
    setIsPlaying(false);
    
    toast({
      title: "进度已重置",
      description: "重新开始背诵练习",
    });
  };

  const getProgress = () => {
    return sentences.length > 0 ? (completedSentences.size / sentences.length) * 100 : 0;
  };

  if (!essay || sentences.length === 0) {
    return (
      <div className="min-h-screen bg-background py-8">
        <div className="container max-w-4xl mx-auto px-4">
          <div className="text-center py-16">
            <h1 className="text-2xl font-bold mb-4">文章加载中...</h1>
            <Button variant="outline" asChild>
              <Link to="/library">
                <ArrowLeft className="w-4 h-4 mr-2" />
                返回文章库
              </Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const currentSentence = sentences[currentSentenceIndex];
  const isCurrentCompleted = completedSentences.has(currentSentenceIndex);

  return (
    <div className="min-h-screen bg-background py-8">
      <div className="container max-w-4xl mx-auto px-4">
        {/* Navigation */}
        <div className="flex items-center justify-between mb-6">
          <Button variant="ghost" asChild>
            <Link to={`/essay/${id}`}>
              <ArrowLeft className="w-4 h-4 mr-2" />
              返回文章详情
            </Link>
          </Button>
          
          <Button variant="outline" onClick={resetProgress}>
            <RotateCcw className="w-4 h-4 mr-2" />
            重置进度
          </Button>
        </div>

        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold mb-2">背诵模式</h1>
          <h2 className="text-lg text-muted-foreground mb-4">{essay.title}</h2>
          
          {/* Progress */}
          <div className="max-w-md mx-auto space-y-2">
            <div className="flex justify-between text-sm">
              <span>进度</span>
              <span>{completedSentences.size} / {sentences.length} 句</span>
            </div>
            <Progress value={getProgress()} className="h-2" />
            <p className="text-xs text-muted-foreground">
              已完成 {Math.round(getProgress())}%
            </p>
          </div>
        </div>

        {/* Main Content */}
        <div className="max-w-2xl mx-auto space-y-6">
          {/* Sentence Display */}
          <Card variant="learning" className="min-h-[200px]">
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg">
                  第 {currentSentenceIndex + 1} / {sentences.length} 句
                </CardTitle>
                <div className="flex items-center gap-2">
                  {isCurrentCompleted && (
                    <Badge variant="secondary" className="bg-learning-success/10 text-learning-success">
                      <CheckCircle className="w-3 h-3 mr-1" />
                      已完成
                    </Badge>
                  )}
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="min-h-[80px] flex items-center justify-center p-4 rounded-lg bg-muted/30">
                {showText ? (
                  <p className="text-lg leading-relaxed text-center">
                    {currentSentence}
                  </p>
                ) : (
                  <div className="text-center space-y-2">
                    <Volume2 className="w-8 h-8 text-muted-foreground mx-auto" />
                    <p className="text-sm text-muted-foreground">
                      点击播放按钮听取音频，然后尝试复述这句话
                    </p>
                  </div>
                )}
              </div>
              
              <div className="flex justify-center">
                <Button variant="outline" onClick={toggleTextVisibility}>
                  {showText ? (
                    <>
                      <EyeOff className="w-4 h-4 mr-2" />
                      隐藏文字
                    </>
                  ) : (
                    <>
                      <Eye className="w-4 h-4 mr-2" />
                      显示文字
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Controls */}
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-center gap-4">
                <Button
                  variant="outline"
                  size="icon"
                  onClick={previousSentence}
                  disabled={currentSentenceIndex === 0}
                >
                  <SkipBack className="w-4 h-4" />
                </Button>
                
                <Button
                  variant="hero"
                  size="lg"
                  onClick={playCurrentSentence}
                  disabled={isPlaying}
                >
                  {isPlaying ? (
                    <Pause className="w-5 h-5" />
                  ) : (
                    <Play className="w-5 h-5" />
                  )}
                </Button>
                
                <Button
                  variant="outline"
                  size="icon"
                  onClick={nextSentence}
                  disabled={currentSentenceIndex === sentences.length - 1}
                >
                  <SkipForward className="w-4 h-4" />
                </Button>
              </div>
              
              <div className="mt-4 text-center">
                <Button
                  variant={isCurrentCompleted ? "outline" : "learning"}
                  onClick={markCurrentAsCompleted}
                  disabled={isCurrentCompleted}
                >
                  {isCurrentCompleted ? (
                    <>
                      <CheckCircle className="w-4 h-4 mr-2" />
                      已掌握
                    </>
                  ) : (
                    <>
                      <CheckCircle className="w-4 h-4 mr-2" />
                      标记为已掌握
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Tips */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">💡 背诵提示</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm text-muted-foreground">
              <p>• 先听音频，理解句子的语音语调</p>
              <p>• 尝试不看文字复述句子</p>
              <p>• 多次练习直到能够流利背诵</p>
              <p>• 完成一句后标记为已掌握</p>
              <p>• 定期回顾已完成的句子</p>
            </CardContent>
          </Card>

          {/* Quick Navigation */}
          {sentences.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">快速导航</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-10 gap-2">
                  {sentences.map((_, index) => (
                    <Button
                      key={index}
                      variant={
                        index === currentSentenceIndex
                          ? "default"
                          : completedSentences.has(index)
                          ? "learning"
                          : "outline"
                      }
                      size="sm"
                      className="aspect-square p-0"
                      onClick={() => {
                        setCurrentSentenceIndex(index);
                        setShowText(false);
                        setIsPlaying(false);
                      }}
                    >
                      {index + 1}
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}