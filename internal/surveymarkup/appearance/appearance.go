package appearance

const (
	Overlay       = "<div class=pure-srv-overlay><div class=pure-srv-float-ctr><div class=pure-srv__close-overlay><svg width=\"40\" height=\"40\" viewbox=\"0 0 40 40\"><path d=\"M10 10 30 30m0-20L10 30\" stroke=\"#000\" stroke-width=\"3\"/></svg></div><div class=pure-srv-float-ctr__flex>${SURVEY_PLACEMENT}</div></div></div><style>.pure-srv-overlay{position:fixed;width:100%;height:100%;top:0;left:0;right:0;bottom:0;background-color:rgba(0,0,0,.5);z-index:999999;cursor:pointer}.pure-srv__close-overlay{position:absolute;top:0;right:0;cursor:pointer;padding:5px}.pure-srv__close-overlay svg{width:32px;height:32px}.pure-srv-float-ctr{width:80%;max-width:640px;height:80%;max-height:500px;position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);cursor:auto;background-color:#fff;padding:30px}.pure-srv-float-ctr__flex{display:flex;justify-content:center;align-items:center;height:100%}</style>"
	OverlayScript = "document.addEventListener(\"keyup\",function(e){e.key===\"Escape\"&&closePureSrvOverlay()}),document.getElementsByClassName(\"pure-srv-overlay\")[0].addEventListener(\"click\",e=>{closePureSrvOverlay()}),document.getElementsByClassName(\"pure-srv__close-overlay\")[0].addEventListener(\"click\",e=>{closePureSrvOverlay()}),document.getElementsByClassName(\"pure-srv-float-ctr\")[0].addEventListener(\"click\",e=>{e.stopPropagation()});function closePureSrvOverlay(){document.getElementsByClassName(\"pure-srv-overlay\")[0].style.display=\"none\"}"
	Left          = "<div class=pure-srv-float-ctr><div class=pure-srv__close-overlay><svg width=\"40\" height=\"40\" viewbox=\"0 0 40 40\"><path d=\"M10 10 30 30m0-20L10 30\" stroke=\"#000\" stroke-width=\"3\"/></svg></div><div class=pure-srv-float-ctr__flex>${SURVEY_PLACEMENT}</div></div><style>.pure-srv__close-overlay{position:absolute;top:0;right:0;cursor:pointer;padding:5px}.pure-srv__close-overlay svg{width:32px;height:32px}.pure-srv-float-ctr{width:250px;max-width:250px;height:auto;position:absolute;top:0;left:0;bottom:0;margin:15px;cursor:auto;background-color:#fff;padding:30px;border:1px solid #000}.pure-srv-float-ctr__flex{display:flex;justify-content:center;align-items:center;height:100%}@media(max-width:550px){.pure-srv-float-ctr{right:0;width:auto;max-width:100%}}</style>"
	LeftScript    = "document.getElementsByClassName(\"pure-srv__close-overlay\")[0].addEventListener(\"click\",e=>{closePureSrvOverlay()});function closePureSrvOverlay(){document.getElementsByClassName(\"pure-srv-float-ctr\")[0].style.display=\"none\"}"
	Right         = "<div class=pure-srv-float-ctr><div class=pure-srv__close-overlay><svg width=\"40\" height=\"40\" viewbox=\"0 0 40 40\"><path d=\"M10 10 30 30m0-20L10 30\" stroke=\"#000\" stroke-width=\"3\"/></svg></div><div class=pure-srv-float-ctr__flex>${SURVEY_PLACEMENT}</div></div><style>.pure-srv__close-overlay{position:absolute;top:0;right:0;cursor:pointer;padding:5px}.pure-srv__close-overlay svg{width:32px;height:32px}.pure-srv-float-ctr{width:250px;max-width:250px;height:auto;position:absolute;top:0;right:0;bottom:0;margin:15px;cursor:auto;background-color:#fff;padding:30px;border:1px solid #000}.pure-srv-float-ctr__flex{display:flex;justify-content:center;align-items:center;height:100%}@media(max-width:550px){.pure-srv-float-ctr{left:0;width:auto;max-width:100%}}</style>"
	RightScript   = "document.getElementsByClassName(\"pure-srv__close-overlay\")[0].addEventListener(\"click\",e=>{closePureSrvOverlay()});function closePureSrvOverlay(){document.getElementsByClassName(\"pure-srv-float-ctr\")[0].style.display=\"none\"}"

	ResponseScript = "const htmlDecode=e=>{var t=document.createElement(\"div\");return t.innerHTML=e,t.childNodes.length===0?\"\":t.childNodes[0].nodeValue};var resp=\"${HTML}\",script,scriptContent=\"${SCRIPT}\",element=document.getElementById(\"pure-srv-ui--${UNIT_ID}\");element.innerHTML=htmlDecode(resp),script=document.createElement(\"script\"),script.innerHTML=htmlDecode(scriptContent),element.appendChild(script)"
)