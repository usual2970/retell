import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Upload, Loader2, CheckCircle, FileText } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { useNavigate } from "react-router-dom";

export default function Import() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim() || !content.trim()) {
      toast({
        title: "请填写完整信息",
        description: "文章标题和内容都不能为空",
        variant: "destructive"
      });
      return;
    }

    setIsLoading(true);
    
    try {
      // 模拟API调用 - 在实际应用中这里会调用真实的API
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // 模拟生成的文章ID
      const essayId = Date.now().toString();
      
      // 保存到localStorage作为示例存储
      const savedEssays = JSON.parse(localStorage.getItem('essays') || '[]');
      const newEssay = {
        id: essayId,
        title,
        content,
        createdAt: new Date().toISOString(),
        audioUrl: "", // 实际应用中会有真实的音频URL
        imageUrl: `https://images.unsplash.com/photo-1649972904349-6e44c42644a7?w=800&h=600&fit=crop&crop=center`,
        wordCount: content.split(/\s+/).filter(word => word.length > 0).length
      };
      
      savedEssays.push(newEssay);
      localStorage.setItem('essays', JSON.stringify(savedEssays));
      
      toast({
        title: "文章导入成功！",
        description: "语音和图片已自动生成，正在跳转到文章详情页...",
      });
      
      // 跳转到文章详情页
      setTimeout(() => {
        navigate(`/essay/${essayId}`);
      }, 1500);
      
    } catch (error) {
      toast({
        title: "导入失败",
        description: "请检查网络连接后重试",
        variant: "destructive"
      });
    } finally {
      setIsLoading(false);
    }
  };

  const exampleEssays = [
    {
      title: "The Benefits of Reading",
      preview: "Reading is one of the most beneficial activities we can engage in. It expands our knowledge, improves our vocabulary..."
    },
    {
      title: "Climate Change Solutions",
      preview: "Climate change is one of the most pressing issues of our time. However, there are many practical solutions..."
    }
  ];

  return (
    <div className="min-h-screen bg-background py-8">
      <div className="container max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-primary mb-4">
            <Upload className="w-6 h-6 text-primary-foreground" />
          </div>
          <h1 className="text-3xl font-bold mb-2">导入新文章</h1>
          <p className="text-muted-foreground text-lg">
            输入你想要背诵的英文文章，我们将自动生成语音和配图
          </p>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Import Form */}
          <div className="lg:col-span-2">
            <Card variant="learning">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="w-5 h-5" />
                  文章信息
                </CardTitle>
                <CardDescription>
                  请填写文章标题和正文内容
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-6">
                  <div className="space-y-2">
                    <Label htmlFor="title">文章标题</Label>
                    <Input
                      id="title"
                      placeholder="例如：The Importance of Education"
                      value={title}
                      onChange={(e) => setTitle(e.target.value)}
                      disabled={isLoading}
                    />
                  </div>
                  
                  <div className="space-y-2">
                    <Label htmlFor="content">文章内容</Label>
                    <Textarea
                      id="content"
                      placeholder="请粘贴或输入英文文章内容..."
                      className="min-h-[300px] resize-none"
                      value={content}
                      onChange={(e) => setContent(e.target.value)}
                      disabled={isLoading}
                    />
                    <div className="text-sm text-muted-foreground">
                      当前字数: {content.split(/\s+/).filter(word => word.length > 0).length} 词
                    </div>
                  </div>
                  
                  <Button 
                    type="submit" 
                    className="w-full" 
                    variant="hero"
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                        正在生成语音和图片...
                      </>
                    ) : (
                      <>
                        <CheckCircle className="w-4 h-4 mr-2" />
                        导入并生成
                      </>
                    )}
                  </Button>
                </form>
              </CardContent>
            </Card>
          </div>

          {/* Tips and Examples */}
          <div className="space-y-6">
            {/* Tips */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">💡 使用提示</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="text-sm space-y-2">
                  <p>• 建议文章长度在100-1000词之间</p>
                  <p>• 请确保文章内容为英文</p>
                  <p>• 系统会自动生成高质量语音</p>
                  <p>• 将根据内容自动配图</p>
                </div>
              </CardContent>
            </Card>

            {/* Example Essays */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">📝 示例文章</CardTitle>
                <CardDescription>点击使用示例文章</CardDescription>
              </CardHeader>
              <CardContent className="space-y-3">
                {exampleEssays.map((essay, index) => (
                  <div
                    key={index}
                    className="p-3 border rounded-lg cursor-pointer hover:bg-muted/50 transition-colors"
                    onClick={() => {
                      setTitle(essay.title);
                      setContent(essay.preview + " (This is a sample essay for demonstration purposes. Please replace with your actual content.)");
                    }}
                  >
                    <h4 className="font-medium text-sm mb-1">{essay.title}</h4>
                    <p className="text-xs text-muted-foreground line-clamp-2">
                      {essay.preview}
                    </p>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}