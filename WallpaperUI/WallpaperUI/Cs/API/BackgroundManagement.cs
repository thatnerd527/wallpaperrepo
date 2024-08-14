using Microsoft.JSInterop;

using WallpaperUI.Pages;

namespace WallpaperUI.Cs.API
{
    public class BackgroundManagement
    {
        [JSInvokable]

        public static void SetBackground(string backgroundid, string backgrounddata)
        {
            var background = Home.availableBackgrounds.Find(x => x.BaseBackground.FixedBackgroundID == backgroundid);
            if (background != null)
            {
                var elems = Home.backgrounds.Where(x => x.BaseBackground.FixedBackgroundID == backgroundid).FirstOrDefault();
                if (elems != null)
                {
                    Home.currentBackground = elems.UniqueBackgroundID;
                } else
                {
                    var instanced = background.Clone();
                    instanced.UniqueBackgroundID = Guid.NewGuid().ToString();
                    Home.backgrounds.Add(instanced);
                    Home.currentBackground = instanced.UniqueBackgroundID;
                }
                App.SaveBackgroundData(Home.backgrounds);
            } else
            {
                return;
            }
        }

        [JSInvokable]
        public static string GetBackground()
        {
            return Home.currentBackground;
        }

        [JSInvokable]
        public static void SetBackgroundData(string persistentbackgroundid, string backgrounddata)
        {
            var background = Home.backgrounds.Find(x => x.UniqueBackgroundID == persistentbackgroundid);
            if (background != null)
            {
                background.PersistentData = backgrounddata;
                App.SaveBackgroundData(Home.backgrounds);
            }

        }

        [JSInvokable]
        public static string GetBackgroundData(string persistentbackgroundid)
        {
            var background = Home.backgrounds.Find(x => x.UniqueBackgroundID == persistentbackgroundid);
            return background == null ? "" : background.PersistentData;
        }
    }
}
