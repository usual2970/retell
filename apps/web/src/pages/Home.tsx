import { Book, Upload, Library, Volume2, Brain, Headphones } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Link } from "react-router-dom";
import heroImage from "@/assets/hero-image.jpg";

export default function Home() {
  const features = [
    {
      icon: Volume2,
      title: "智能语音生成",
      description: "自动将你的文章转换为高质量语音，支持多种语调和语速"
    },
    {
      icon: Brain,
      title: "背诵训练",
      description: "逐句播放，隐藏文字，帮助你更好地记忆和背诵英文文章"
    },
    {
      icon: Headphones,
      title: "听力练习",
      description: "结合音频和文本，提升你的英语听力理解能力"
    }
  ];

  const quickActions = [
    {
      title: "导入新文章",
      description: "开始你的学习之旅",
      href: "/import",
      icon: Upload,
      variant: "hero" as const
    },
    {
      title: "浏览文章库",
      description: "查看你的学习历史",
      href: "/library",
      icon: Library,
      variant: "learning" as const
    }
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="relative py-20 px-4 overflow-hidden">
        <div className="absolute inset-0 bg-gradient-hero opacity-5" />
        <div 
          className="absolute inset-0 bg-cover bg-center opacity-10"
          style={{ backgroundImage: `url(${heroImage})` }}
        />
        <div className="container max-w-4xl mx-auto text-center relative">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-primary mb-6">
            <Book className="w-8 h-8 text-primary-foreground" />
          </div>
          
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-primary bg-clip-text text-transparent">
            Essay Audio Genie
          </h1>
          
          <p className="text-xl md:text-2xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            让英文essay背诵变得简单有趣。自动生成语音，智能辅助练习，提升你的英语学习效果。
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" variant="hero" asChild>
              <Link to="/import">
                <Upload className="w-5 h-5 mr-2" />
                开始导入文章
              </Link>
            </Button>
            
            <Button size="lg" variant="outline" asChild>
              <Link to="/library">
                <Library className="w-5 h-5 mr-2" />
                查看文章库
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 px-4 bg-muted/30">
        <div className="container max-w-6xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">核心功能特色</h2>
          
          <div className="grid md:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <Card key={index} variant="learning" className="text-center">
                <CardHeader>
                  <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-primary mx-auto mb-4">
                    <feature.icon className="w-6 h-6 text-primary-foreground" />
                  </div>
                  <CardTitle className="text-xl">{feature.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-base">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Quick Actions */}
      <section className="py-16 px-4">
        <div className="container max-w-4xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">快速开始</h2>
          
          <div className="grid md:grid-cols-2 gap-8">
            {quickActions.map((action, index) => (
              <Card key={index} variant="learning" className="hover:scale-105 transition-transform duration-300">
                <CardHeader className="text-center">
                  <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-accent mx-auto mb-4">
                    <action.icon className="w-6 h-6 text-accent-foreground" />
                  </div>
                  <CardTitle className="text-xl">{action.title}</CardTitle>
                  <CardDescription>{action.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <Button className="w-full" variant={action.variant} asChild>
                    <Link to={action.href}>
                      开始使用
                    </Link>
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-16 px-4 bg-muted/30">
        <div className="container max-w-4xl mx-auto text-center">
          <h2 className="text-3xl font-bold mb-12">使用步骤</h2>
          
          <div className="grid md:grid-cols-3 gap-8">
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-gradient-primary text-primary-foreground flex items-center justify-center text-xl font-bold mb-4">
                1
              </div>
              <h3 className="text-lg font-semibold mb-2">导入文章</h3>
              <p className="text-muted-foreground">粘贴或输入你想要背诵的英文文章</p>
            </div>
            
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-gradient-primary text-primary-foreground flex items-center justify-center text-xl font-bold mb-4">
                2
              </div>
              <h3 className="text-lg font-semibold mb-2">自动生成</h3>
              <p className="text-muted-foreground">系统自动生成语音和配图</p>
            </div>
            
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-gradient-primary text-primary-foreground flex items-center justify-center text-xl font-bold mb-4">
                3
              </div>
              <h3 className="text-lg font-semibold mb-2">开始练习</h3>
              <p className="text-muted-foreground">使用背诵模式练习记忆</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}