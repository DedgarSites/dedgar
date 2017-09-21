$(function(){
    
    var sizeFlg;
    var win=$(window).width();
    var p=638;
    

    if(win>p){
        $("ul#sideNavi").show();
        sizeFlg="pc";
        
    }else{
        $("ul#sideNavi").hide();
        sizeFlg="smp";
    }
    
    //#backTopのクリック処理
	$('#backTop').on("click",function(){
		$('body,html').animate({scrollTop:0},400,'swing');
		return false;
									  
	});

    /*スマホメニュー
    ----------------------------------------------------------------------------------------*/
    $("#smpArea .smpBt").click(function(){
        
        $("ul#sideNavi").stop().slideToggle();
        return false;
    });
    /*リサイズ
    ----------------------------------------------------------------------------------------*/
    $(window).resize(function(){
        win=$(window).width();
        
        if(win>p){
            $("ul#sideNavi").show();
            sizeFlg="pc";
            
        }else{
            $("ul#sideNavi").hide();
            sizeFlg="smp";
        }
        xResize();
    });
    /*クリック挙動
    ----------------------------------------------------------------------------------------*/
    $("a[href^=#]").click(function(){
       
		var speed = 400;
		var href = $(this).attr("href");
		var target = $(href == "#" || href == "" ? 'html' : href);
        if(sizeFlg=="pc"){
            var position = Math.ceil(target.offset().top);
        }else if(sizeFlg=="smp"){
		    var position = Math.ceil(target.offset().top-40);
            $("ul#sideNavi").hide();
        }
        
		$('body,html').animate({scrollTop:position},speed,'swing');
		return false;
    });
    
        
    /*スクロール関係
    ----------------------------------------------------------------------------------------*/
    if(sizeFlg=="pc"){
        var mapHight = Math.floor($('#map').height());
        var introPos = Math.floor($('#intro').offset().top);
        var entryPos = Math.floor($('#entry').offset().top);
        var job1Pos = Math.floor($('#job1').offset().top);//プロデューサー/ディレクター
        var job2Pos = Math.floor($('#job2').offset().top);//プランナー
        var job3Pos = Math.floor($('#job3').offset().top);//シナリオライター
        var job4Pos = Math.floor($('#job4').offset().top);//ゲームプログラマー
        var job5Pos = Math.floor($('#job5').offset().top);//アートデザイナー
        var job6Pos = Math.floor($('#job6').offset().top);//３Dデザイナー
        var job7Pos = Math.floor($('#job7').offset().top);//ムービーデザイナー
        var job8Pos = Math.floor($('#job8').offset().top);//Webデザイナー
        
    }else if(sizeFlg=="smp"){
        var mapHight = Math.floor($('#map').height()-40);
        var introPos = Math.floor($('#intro').offset().top-40);
        var entryPos = Math.floor($('#entry').offset().top-40);
        var job1Pos = Math.floor($('#job1').offset().top-40);
        var job2Pos = Math.floor($('#job2').offset().top-40);
        var job3Pos = Math.floor($('#job3').offset().top-40);
        var job4Pos = Math.floor($('#job4').offset().top-40);
        var job5Pos = Math.floor($('#job5').offset().top-40);
        var job6Pos = Math.floor($('#job6').offset().top-40);
        var job7Pos = Math.floor($('#job7').offset().top-40);
        var job8Pos = Math.floor($('#job8').offset().top-40);
    }

    $(window).scroll(function(){
        //TOPへもどる
        if($(window).scrollTop()>100){
			  $('#backTop').stop().animate({"opacity":1},1000,'easeOutExpo');
		}else{
			 
			  $('#backTop').stop().animate({"opacity":0},1000,'easeOutExpo');
		}
        //-----------------------------------------------------------

        
            var y = $(this).scrollTop();		
            $('#scrollValue').text(y);
            if(y<mapHight-10){
                 $('ul#sideNavi li a').removeClass("selected");
            }
            
            if(sizeFlg=="pc"){
                if(y>=mapHight-20){
                     $('ul#sideNavi').stop().animate({"right":"0px"},1000,'easeOutExpo');
                     
                }else{
                     
                      $('ul#sideNavi').stop().animate({"right":"-100px"},1000,'easeOutExpo');
                }
            }
       
            if(y>=mapHight){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(0)').addClass("selected");
            }
           if(y>=entryPos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(1)').addClass("selected");
            }
             
            if(y>=job1Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(2)').addClass("selected");
            }
            if(y>=job2Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(3)').addClass("selected");
            }
            if(y>=job3Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(4)').addClass("selected");
            }
            if(y>=job4Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(5)').addClass("selected");
            }
            if(y>=job5Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(6)').addClass("selected");
            }
            if(y>=job6Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(7)').addClass("selected");
            }
            if(y>=job7Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(8)').addClass("selected");
            }
            if(y>=job8Pos){
                $('ul#sideNavi li a').removeClass("selected");
                $('ul#sideNavi li a:eq(9)').addClass("selected");
            }
    });
    $("#wrap").stop().css({opacity:0});
    cEffectEnd();
});


//リサイズ処理まとめ
function xResize(){
    $(function(){
    if(sizeFlg=="pc"){
        var mapHight = Math.floor($('#map').height());
        var introPos = Math.floor($('#intro').offset().top);
        var entryPos = Math.floor($('#entry').offset().top);
        var job1Pos = Math.floor($('#job1').offset().top);
        var job2Pos = Math.floor($('#job2').offset().top);
        var job3Pos = Math.floor($('#job3').offset().top);
        var job4Pos = Math.floor($('#job4').offset().top);
        var job5Pos = Math.floor($('#job5').offset().top);
        var job6Pos = Math.floor($('#job6').offset().top);
        var job7Pos = Math.floor($('#job7').offset().top);
        var job8Pos = Math.floor($('#job8').offset().top);
        
    }else if(sizeFlg=="smp"){
        var mapHight = Math.floor($('#map').height()-40);
        var introPos = Math.floor($('#intro').offset().top-40);
        var entryPos = Math.floor($('#entry').offset().top-40);
        var job1Pos = Math.floor($('#job1').offset().top-40);
        var job2Pos = Math.floor($('#job2').offset().top-40);
        var job3Pos = Math.floor($('#job3').offset().top-40);
        var job4Pos = Math.floor($('#job4').offset().top-40);
        var job5Pos = Math.floor($('#job5').offset().top-40);
        var job6Pos = Math.floor($('#job6').offset().top-40);
        var job7Pos = Math.floor($('#job7').offset().top-40);
        var job8Pos = Math.floor($('#job8').offset().top-40);
    
    }
    });
}

function cEffectEnd(){
	
	//配列に使用画像URLを代入------------------------
	var pre = [];
	
	//div画像背景
	for (var i =0; i < $("div").length; i++){
		
		if($("div:eq("+i+")").css("background-image") != "none"){
		
			var preImg = $("div:eq("+i+")").css("background-image");
		
			if (preImg.indexOf("gradient") == -1){
				preImg = preImg.replace(/url\(\"/g, "");
				preImg = preImg.replace(/url\(/g, "");
				preImg = preImg.replace(/\"\)/g, "");
				preImg = preImg.replace(/\)/g, "");
			}
			pre.push(preImg);
		}
	}
	//article画像背景
	for (var i =0; i < $("article").length; i++){
		
		if($("article:eq("+i+")").css("background-image") != "none"){
		
			var preImg = $("article:eq("+i+")").css("background-image");
		
			if (preImg.indexOf("gradient") == -1){
				preImg = preImg.replace(/url\(\"/g, "");
				preImg = preImg.replace(/url\(/g, "");
				preImg = preImg.replace(/\"\)/g, "");
				preImg = preImg.replace(/\)/g, "");
			}
			pre.push(preImg);
		}
	}
    //section画像背景
	for (var i =0; i < $("section").length; i++){
		
		if($("section:eq("+i+")").css("background-image") != "none"){
		
			var preImg = $("section:eq("+i+")").css("background-image");
		
			if (preImg.indexOf("gradient") == -1){
				preImg = preImg.replace(/url\(\"/g, "");
				preImg = preImg.replace(/url\(/g, "");
				preImg = preImg.replace(/\"\)/g, "");
				preImg = preImg.replace(/\)/g, "");
			}
			pre.push(preImg);
		}
	}

	//img画像
	for (var i =0; i < $("img").length; i++){
		var preImg = $("img:eq("+i+")").attr("src");
		pre.push(preImg);
	}
	//alert(pre)
	//-------------------------------------------------------
	//オブジェクトがフェードアウトするまで待つ
	
	setTimeout(function(){
		 
		//配列削除用
		Array.prototype.remove = function(element){
		  for(var i =0; i < this.length; i++){
			if (this[i] == element) this.splice(i,1);
		  }
		};

		function preload(img, progress) {
			var total = img.length;

			$(img).each(function(){
				var src = this;
				$('<img />').on("load",function(){
					img.remove(src);
					progress(total, total - img.length);
				}).attr('src',src);
				//《.on("load"》してからの《.attr》が重要なので順番は変えないようお願いします。
			});
		}
		
		//ロード実行
		preload(pre,function(total, loaded){
           
		  if(loaded >= total){
               
				 $("#wrap").stop().css({opacity:0}).show().animate({opacity:1},1000,'easeOutExpo');
		  }else{
          }
		});

	},500);
			
};