using Microsoft.JSInterop;

using WallpaperUI.Pages;

namespace WallpaperUI.Cs.API
{
    public class BackgroundManagement
    {
        [JSInvokable]

        public static void SetBackground(string backgroundid, string backgrounddata)
        {
            var background = Home.availableBackgrounds.Find(x => x.LoaderBackgroundID == backgroundid);
            if (background != null)
            {
                var elems = Home.backgrounds.Where(x => x.LoaderBackgroundID == backgroundid).FirstOrDefault();
                if (elems != null)
                {
                    Home.currentBackground = elems.PersistentBackgroundID;
                } else
                {
                    var instanced = background.Cloned();
                    instanced.PersistentBackgroundID = Guid.NewGuid().ToString();
                    Home.backgrounds.Add(instanced);
                    Home.currentBackground = instanced.PersistentBackgroundID;
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
            var background = Home.backgrounds.Find(x => x.PersistentBackgroundID == persistentbackgroundid);
            if (background != null)
            {
                background.PersistentBackgroundData = backgrounddata;
                App.SaveBackgroundData(Home.backgrounds);
            }

        }

        [JSInvokable]
        public static string GetBackgroundData(string persistentbackgroundid)
        {
            var background = Home.backgrounds.Find(x => x.PersistentBackgroundID == persistentbackgroundid);
            return background == null ? "" : background.PersistentBackgroundData;
        }
    }
}
