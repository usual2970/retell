import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Library, Search, Clock, Volume2, Eye, Plus, BookOpen } from "lucide-react";
import { Link } from "react-router-dom";
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

export default function LibraryPage() {
  const [essays, setEssays] = useState<Essay[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [filteredEssays, setFilteredEssays] = useState<Essay[]>([]);

  useEffect(() => {
    // 从localStorage加载文章
    const savedEssays = JSON.parse(localStorage.getItem('essays') || '[]');
    setEssays(savedEssays);
    setFilteredEssays(savedEssays);
  }, []);

  useEffect(() => {
    // 搜索过滤
    if (searchQuery.trim() === "") {
      setFilteredEssays(essays);
    } else {
      const filtered = essays.filter(essay =>
        essay.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        essay.content.toLowerCase().includes(searchQuery.toLowerCase())
      );
      setFilteredEssays(filtered);
    }
  }, [searchQuery, essays]);

  const getReadingTime = (wordCount: number) => {
    const wordsPerMinute = 200;
    const minutes = Math.ceil(wordCount / wordsPerMinute);
    return `${minutes} 分钟`;
  };

  const truncateContent = (content: string, maxLength: number = 150) => {
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + "...";
  };

  if (essays.length === 0) {
    return (
      <div className="min-h-screen bg-background py-8">
        <div className="container max-w-4xl mx-auto px-4">
          <div className="text-center py-16">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-6">
              <Library className="w-8 h-8 text-muted-foreground" />
            </div>
            <h1 className="text-2xl font-bold mb-4">文章库为空</h1>
            <p className="text-muted-foreground mb-8">
              你还没有导入任何文章。立即开始导入你的第一篇文章吧！
            </p>
            <Button variant="hero" size="lg" asChild>
              <Link to="/import">
                <Plus className="w-5 h-5 mr-2" />
                导入第一篇文章
              </Link>
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background py-8">
      <div className="container max-w-6xl mx-auto px-4">
        {/* Header */}
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold mb-2 flex items-center gap-3">
              <Library className="w-8 h-8 text-primary" />
              我的文章库
            </h1>
            <p className="text-muted-foreground">
              共 {essays.length} 篇文章 • {filteredEssays.length} 篇符合条件
            </p>
          </div>
          
          <Button variant="hero" asChild>
            <Link to="/import">
              <Plus className="w-4 h-4 mr-2" />
              导入新文章
            </Link>
          </Button>
        </div>

        {/* Search */}
        <div className="mb-8">
          <div className="relative max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
            <Input
              placeholder="搜索文章标题或内容..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {/* Essays Grid */}
        {filteredEssays.length === 0 ? (
          <div className="text-center py-12">
            <Search className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">未找到匹配的文章</h3>
            <p className="text-muted-foreground">
              尝试使用不同的关键词搜索
            </p>
          </div>
        ) : (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredEssays.map((essay) => (
              <Card key={essay.id} variant="learning" className="group">
                <CardHeader className="pb-3">
                  <div className="aspect-video rounded-lg overflow-hidden mb-3">
                    <img
                      src={essay.imageUrl}
                      alt={essay.title}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                  </div>
                  <CardTitle className="text-lg line-clamp-2">
                    {essay.title}
                  </CardTitle>
                  <CardDescription className="line-clamp-3">
                    {truncateContent(essay.content)}
                  </CardDescription>
                </CardHeader>
                
                <CardContent className="space-y-4">
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Clock className="w-4 h-4" />
                      {getReadingTime(essay.wordCount)}
                    </div>
                    <Badge variant="secondary">
                      {essay.wordCount} 词
                    </Badge>
                  </div>
                  
                  <div className="text-xs text-muted-foreground">
                    创建时间：{format(new Date(essay.createdAt), 'yyyy年MM月dd日', { locale: zhCN })}
                  </div>
                  
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm" asChild className="flex-1">
                      <Link to={`/essay/${essay.id}`}>
                        <Eye className="w-4 h-4 mr-1" />
                        查看
                      </Link>
                    </Button>
                    
                    <Button variant="accent" size="sm" asChild className="flex-1">
                      <Link to={`/essay/${essay.id}/recite`}>
                        <BookOpen className="w-4 h-4 mr-1" />
                        背诵
                      </Link>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Stats */}
        <div className="mt-12 grid md:grid-cols-3 gap-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-primary/10">
                  <Library className="w-5 h-5 text-primary" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{essays.length}</p>
                  <p className="text-sm text-muted-foreground">总文章数</p>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-accent/10">
                  <Volume2 className="w-5 h-5 text-accent" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {essays.reduce((total, essay) => total + essay.wordCount, 0)}
                  </p>
                  <p className="text-sm text-muted-foreground">总词汇量</p>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-learning-success/10">
                  <Clock className="w-5 h-5 text-learning-success" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {Math.ceil(essays.reduce((total, essay) => total + essay.wordCount, 0) / 200)}
                  </p>
                  <p className="text-sm text-muted-foreground">预计阅读时间（分钟）</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}