import { Moon, Sun } from "lucide-react";
import { Switch } from "./switch";
import { useEffect, useState } from "react";

type Theme = 'light' | 'dark';

const ThemeSwitcher = () =>{
    const [theme,setTheme] = useState<Theme>('light');

    useEffect(()=>{
        localStorage.setItem('theme',theme);
        if(theme=='dark'){
            document.documentElement.classList.add('dark');
        }else{
            document.documentElement.classList.remove('dark');
        }
    },[theme])

    const themeToggle = ()=>{
        setTheme(theme === 'dark' ? 'light' : 'dark');
    }

    return(
        <div className="flex gap-2 items-center" >
            <Sun
            className={`h-5 w-5  ${theme==='dark' ? 'text-primary/50' : 'text-primary'}`}/>
            <Switch
            className="border border-gray-400"
            checked={theme === "dark"}
            onCheckedChange={themeToggle}
            />
            <Moon
            className={`h-5 w-5 ${theme==='dark' ? 'text-primary' : 'text-primary/50'}`}/>

        </div>
    )
}
export default ThemeSwitcher;
