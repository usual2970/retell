import { useState, useEffect, useRef } from "react";
import { useParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Play, Pause, RotateCcw, Volume2, BookOpen, ArrowLeft, Clock, FileText, Download } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { format } from "date-fns";
import { zhCN } from "date-fns/locale";

interface Essay {
  id: string;
  title: string;
  content: string;
  createdAt: string;
  audioUrl: string;
  imageUrl: string;
  wordCount: number;
}

export default function EssayDetail() {
  const { id } = useParams<{ id: string }>();
  const [essay, setEssay] = useState<Essay | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [playbackRate, setPlaybackRate] = useState(1);
  const audioRef = useRef<HTMLAudioElement>(null);
  const { toast } = useToast();

  useEffect(() => {
    if (id) {
      const savedEssays = JSON.parse(localStorage.getItem('essays') || '[]');
      const foundEssay = savedEssays.find((e: Essay) => e.id === id);
      if (foundEssay) {
        setEssay(foundEssay);
      }
    }
  }, [id]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const updateTime = () => setCurrentTime(audio.currentTime);
    const updateDuration = () => setDuration(audio.duration);
    const handleEnded = () => setIsPlaying(false);

    audio.addEventListener('timeupdate', updateTime);
    audio.addEventListener('loadedmetadata', updateDuration);
    audio.addEventListener('ended', handleEnded);

    return () => {
      audio.removeEventListener('timeupdate', updateTime);
      audio.removeEventListener('loadedmetadata', updateDuration);
      audio.removeEventListener('ended', handleEnded);
    };
  }, [essay]);

  const togglePlayPause = () => {
    const audio = audioRef.current;
    if (!audio) return;

    if (isPlaying) {
      audio.pause();
    } else {
      // 由于这是演示，我们模拟音频播放
      toast({
        title: "音频功能演示",
        description: "在实际应用中，这里会播放生成的语音",
      });
      // 模拟播放状态
      setIsPlaying(true);
      setTimeout(() => setIsPlaying(false), 3000);
    }
    setIsPlaying(!isPlaying);
  };

  const resetAudio = () => {
    const audio = audioRef.current;
    if (audio) {
      audio.currentTime = 0;
      setCurrentTime(0);
    }
    setIsPlaying(false);
  };

  const changePlaybackRate = () => {
    const rates = [0.5, 0.75, 1, 1.25, 1.5, 2];
    const currentIndex = rates.indexOf(playbackRate);
    const nextRate = rates[(currentIndex + 1) % rates.length];
    setPlaybackRate(nextRate);
    
    const audio = audioRef.current;
    if (audio) {
      audio.playbackRate = nextRate;
    }
  };

  const formatTime = (time: number) => {
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const getReadingTime = (wordCount: number) => {
    return Math.ceil(wordCount / 200);
  };

  if (!essay) {
    return (
      <div className="min-h-screen bg-background py-8">
        <div className="container max-w-4xl mx-auto px-4">
          <div className="text-center py-16">
            <h1 className="text-2xl font-bold mb-4">文章未找到</h1>
            <p className="text-muted-foreground mb-8">
              请检查链接是否正确，或返回文章库重新选择
            </p>
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

  return (
    <div className="min-h-screen bg-background py-8">
      <div className="container max-w-4xl mx-auto px-4">
        {/* Navigation */}
        <div className="mb-6">
          <Button variant="ghost" asChild>
            <Link to="/library">
              <ArrowLeft className="w-4 h-4 mr-2" />
              返回文章库
            </Link>
          </Button>
        </div>

        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-4">{essay.title}</h1>
          <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <div className="flex items-center gap-1">
              <Clock className="w-4 h-4" />
              创建于 {format(new Date(essay.createdAt), 'yyyy年MM月dd日 HH:mm', { locale: zhCN })}
            </div>
            <Badge variant="secondary">{essay.wordCount} 词</Badge>
            <Badge variant="outline">{getReadingTime(essay.wordCount)} 分钟阅读</Badge>
          </div>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Cover Image */}
            <Card>
              <CardContent className="p-0">
                <div className="aspect-video rounded-lg overflow-hidden">
                  <img
                    src={essay.imageUrl}
                    alt={essay.title}
                    className="w-full h-full object-cover"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Article Content */}
            <Card variant="learning">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="w-5 h-5" />
                  文章内容
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="prose prose-gray max-w-none">
                  <div className="whitespace-pre-wrap text-base leading-relaxed">
                    {essay.content}
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Audio Player */}
            <Card variant="learning">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Volume2 className="w-5 h-5" />
                  音频播放器
                </CardTitle>
                <CardDescription>
                  AI生成的高质量语音
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Audio Element (hidden) */}
                <audio ref={audioRef} preload="metadata">
                  <source src={essay.audioUrl} type="audio/mpeg" />
                  您的浏览器不支持音频播放
                </audio>

                {/* Progress Bar */}
                <div className="space-y-2">
                  <div className="w-full bg-muted rounded-full h-2">
                    <div
                      className="bg-primary h-2 rounded-full transition-all duration-300"
                      style={{
                        width: duration > 0 ? `${(currentTime / duration) * 100}%` : '0%'
                      }}
                    />
                  </div>
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>{formatTime(currentTime)}</span>
                    <span>{formatTime(duration)}</span>
                  </div>
                </div>

                {/* Controls */}
                <div className="flex items-center justify-center gap-3">
                  <Button variant="outline" size="icon" onClick={resetAudio}>
                    <RotateCcw className="w-4 h-4" />
                  </Button>
                  
                  <Button variant="hero" size="icon" onClick={togglePlayPause}>
                    {isPlaying ? (
                      <Pause className="w-5 h-5" />
                    ) : (
                      <Play className="w-5 h-5" />
                    )}
                  </Button>
                  
                  <Button variant="outline" size="sm" onClick={changePlaybackRate}>
                    {playbackRate}x
                  </Button>
                </div>

                <p className="text-xs text-muted-foreground text-center">
                  *这是演示版本，实际应用中会有真实的音频文件
                </p>
              </CardContent>
            </Card>

            {/* Actions */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">学习模式</CardTitle>
                <CardDescription>
                  选择你的学习方式
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button className="w-full" variant="accent" asChild>
                  <Link to={`/essay/${essay.id}/recite`}>
                    <BookOpen className="w-4 h-4 mr-2" />
                    开始背诵模式
                  </Link>
                </Button>
                
                <Button className="w-full" variant="outline">
                  <Download className="w-4 h-4 mr-2" />
                  下载音频文件
                </Button>
              </CardContent>
            </Card>

            {/* Statistics */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">文章统计</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-sm text-muted-foreground">总词数</span>
                  <span className="font-medium">{essay.wordCount}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm text-muted-foreground">估计阅读时间</span>
                  <span className="font-medium">{getReadingTime(essay.wordCount)} 分钟</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm text-muted-foreground">创建时间</span>
                  <span className="font-medium">
                    {format(new Date(essay.createdAt), 'MM-dd', { locale: zhCN })}
                  </span>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}