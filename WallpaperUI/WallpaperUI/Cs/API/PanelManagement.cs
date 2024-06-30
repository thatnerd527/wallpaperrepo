using Microsoft.JSInterop;

using WallpaperUI.Pages;

namespace WallpaperUI.Cs.API
{
    public class PanelManagement
    {
        [JSInvokable]
        public static void ClosePanel(string persistentpanelid)
        {
            Home.panels.RemoveAll(x => x.PersistentPanelID == persistentpanelid);
            Home.Instance.UpdatePanelData();
        }

        [JSInvokable]
        public static void OpenPanel(string panelid) {
            var panel = Home.Instance.AddPanel(panelid);
            App.SavePanelData(Home.panels);
            Home.Instance.UpdatePanelData();
        }

        [JSInvokable]
        public static string GetPanelSize(string persistentpanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            return panel == null ? "" : $"{panel.PersistentPanelWidth},{panel.PersistentPanelHeight}";
        }

        [JSInvokable]
        public static void SetPanelSize(string persistentpanelid, double width, double height)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            if (panel != null)
            {
                panel.PersistentPanelWidth = width;
                panel.PersistentPanelHeight = height;
                Home.Instance.UpdatePanelData();
            }
            App.SavePanelData(Home.panels);
        }

        [JSInvokable]
        public static string GetPanelPosition(string persistentpanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            return panel == null ? "" : $"{panel.PersistentPanelX},{panel.PersistentPanelY}";
        }

        [JSInvokable]
        public static void SetPanelPosition(string persistentpanelid, double x, double y)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            if (panel != null)
            {
                panel.PersistentPanelX = x;
                panel.PersistentPanelY = y;
                Home.Instance.UpdatePanelData();
            }
            App.SavePanelData(Home.panels);
        }

        [JSInvokable]
        public static string GetPanelVisibility(string persistentpanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            return panel == null ? "" : panel.Deleted ? "hidden" : "visible";
        }

        [JSInvokable]
        public static void SetPanelData(string persistentpanelid, string data)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID==persistentpanelid);
            if (panel != null)
            {
                panel.PersistentPanelData = data;
                Home.Instance.UpdatePanelData();
            }
        }

        [JSInvokable]
        public static string GetPanelData(string persistentpanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            return panel == null ? "" : panel.PersistentPanelData;
        }

        [JSInvokable]
        public static void SetPanelHeader(string persistentpanelid, bool titlebarvisible, bool enableresize, bool enableclose, bool enabledrag)
        {
            var panel = Home.panels.FirstOrDefault(x => x.PersistentPanelID == persistentpanelid);
            var exists = Home.panelHeaders.Any(x => x.persistentpanelid == persistentpanelid);
            if (exists)
            {
                var res = Home.panelHeaders.FirstOrDefault(x => x.persistentpanelid==persistentpanelid)!;
                res.titlebarvisible = titlebarvisible;
                res.enableresize = enableresize;
                res.enableclose = enableclose;
                res.enabledrag = enabledrag;
            } else
            {
                var header = new PanelHeader();
                header.persistentpanelid = persistentpanelid;
                header.titlebarvisible = titlebarvisible;
                header.enableresize = enableresize;
                header.enableclose = enableclose;
                header.enabledrag = enabledrag;
                Home.panelHeaders.Add(header);
            }
        }

        [JSInvokable]
        public static string GetPanelHeader(string persistentpanelid)
        {
            var header = Home.panelHeaders.FirstOrDefault(x => x.persistentpanelid == persistentpanelid);
            return header == null ? "" : $"{header.titlebarvisible},{header.enableresize},{header.enableclose},{header.enabledrag}";
        }


    }
}
