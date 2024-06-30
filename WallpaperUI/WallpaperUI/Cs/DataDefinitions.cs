﻿using System.Text.Json;
using System.Text.Json.Serialization;

namespace WallpaperUI.Cs
{
    public class TemplateCustomBackground
    {
        public string LoaderBackgroundID { get; }
        public string BackgroundType { get; }
        public string BackgroundContent { get; }
        public string BackgroundDefaultData { get; }
        public string ClientID { get; }

        [JsonConstructor]
        public TemplateCustomBackground(string LoaderBackgroundID, string BackgroundType, string BackgroundContent, string BackgroundDefaultData, string ClientID)
        {
            this.LoaderBackgroundID = LoaderBackgroundID;
            this.BackgroundType = BackgroundType;
            this.BackgroundContent = BackgroundContent;
            this.BackgroundDefaultData = BackgroundDefaultData;
            this.ClientID = ClientID;
        }

        public static TemplateCustomBackground FromJsonElement(JsonElement x)
        {
            TemplateCustomBackground tcp = new TemplateCustomBackground(
                x.GetProperty("LoaderBackgroundID").GetString()!,
                x.GetProperty("BackgroundType").GetString()!,
                x.GetProperty("BackgroundContent").GetString()!,
                x.GetProperty("BackgroundDefaultData").GetString()!,
                x.GetProperty("ClientID").GetString()!
                );
            return tcp;

        }

        public TemplateCustomBackground Cloned()
        {
            var serialized = JsonSerializer.Serialize( this );
            return JsonSerializer.Deserialize<TemplateCustomBackground>(serialized);
        }
    }

    public class RuntimeCustomBackground : TemplateCustomBackground
    {
        public string PersistentBackgroundID { get; set; }
        public string PersistentBackgroundData { get; set; }
        public bool Deleted { get; }
        public int ControlPort { get; }

        [JsonConstructor]
        public RuntimeCustomBackground(string LoaderBackgroundID, string BackgroundType, string BackgroundContent, string BackgroundDefaultData, string ClientID, string PersistentBackgroundID, string PersistedBackgroundData, bool Deleted, int ControlPort) : base(LoaderBackgroundID, BackgroundType, BackgroundContent, BackgroundDefaultData, ClientID)
        {
            this.PersistentBackgroundID = PersistentBackgroundID;
            this.PersistentBackgroundData = PersistedBackgroundData;
            this.Deleted = Deleted;
            this.ControlPort = ControlPort;
        }

        public static RuntimeCustomBackground FromJsonElement(JsonElement element)
        {
            RuntimeCustomBackground rcp = new RuntimeCustomBackground(
                    element.GetProperty("LoaderBackgroundID").GetString()!,
                    element.GetProperty("BackgroundType").GetString()!,
                    element.GetProperty("BackgroundContent").GetString()!,
                    element.GetProperty("BackgroundDefaultData").GetString()!,
                    element.GetProperty("ClientID").GetString()!,
                    element.GetProperty("PersistentBackgroundID").GetString()!,
                    element.GetProperty("PersistentBackgroundData").GetString()!,
                    element.GetProperty("Deleted").GetBoolean(),
                    element.GetProperty("ControlPort").GetInt32()
                );
            return rcp;
        }

        public RuntimeCustomBackground Cloned()
        {
            return new RuntimeCustomBackground(LoaderBackgroundID, BackgroundType, BackgroundContent, BackgroundDefaultData, ClientID
                ,PersistentBackgroundID, PersistentBackgroundData, Deleted, ControlPort);
        }
    }

    public class BackgroundUpdate
    {

        public List<RuntimeCustomBackground> NewBackgrounds { get; set; }
        public string NewActiveBackground { get; set; }
    }

    public class AddonManifest
    {
        public string Name { get; set; }
        public string Version { get; set; }
        public string Description { get; set; }
        public string Author { get; set; }
        public string ClientID { get; set; }
        public string BootstrapExecutable { get; set; }
        public bool EnableRestart { get; set; }

        public AddonManifest(string name, string version, string description, string author, string clientID, string bootstrapExecutable, bool enableRestart)
        {
            Name = name;
            Version = version;
            Description = description;
            Author = author;
            ClientID = clientID;
            BootstrapExecutable = bootstrapExecutable;
            EnableRestart = enableRestart;
        }

        public static AddonManifest FromJsonElement(JsonElement x)
        {
            AddonManifest addonManifest = new AddonManifest(
                x.GetProperty("name").GetString()!,
                x.GetProperty("version").GetString()!,
                x.GetProperty("description").GetString()!,
                x.GetProperty("author").GetString()!,
                x.GetProperty("clientID").GetString()!,
                x.GetProperty("bootstrapExecutable").GetString()!,
                x.GetProperty("enableRestart").GetBoolean()
                );
            return addonManifest;
        }

        public AddonManifest Cloned()
        {
            var serialized = JsonSerializer.Serialize( this );
            return JsonSerializer.Deserialize<AddonManifest>(serialized);
        }
    }



    public class RuntimeCustomPanel
    {
        public string PanelType { get; }
        public string LoaderPanelID { get; }
        public string PersistentPanelID { get; set; }
        public string PanelTitle { get; }
        public string PanelContent { get; }
        public double PanelRecommendedWidth { get; }
        public double PanelRecommendedHeight { get; }
        public double PanelMinWidth { get; }
        public double PanelMinHeight { get; }
        public double PanelMaxWidth { get; }
        public double PanelMaxHeight { get; }
        public double PersistentPanelWidth { get; set; }
        public double PersistentPanelHeight { get; set; }
        public double PanelRecommendedX { get; }
        public double PanelRecommendedY { get; }
        public double PersistentPanelX { get; set; }
        public double PersistentPanelY { get; set; }
        public int ControlPort { get; }
        public bool Deleted { get; }
        public string PersistentPanelData { get; set; }
        public string PanelIcon { get; }
        public string ClientID { get; }

        public RuntimeCustomPanel(string panelType, string loaderPanelID, string persistentPanelID, string panelTitle, string panelContent, double panelRecommendedWidth, double panelRecommendedHeight, double panelMinWidth, double panelMinHeight, double panelMaxWidth, double panelMaxHeight, double persistentPanelWidth, double persistentPanelHeight, double panelRecommendedX, double panelRecommendedY, double persistentPanelX, double persistentPanelY, int controlPort, bool deleted, string persistentPanelData, string panelIcon, string clientID)
        {
            PanelType = panelType;
            LoaderPanelID = loaderPanelID;
            PersistentPanelID = persistentPanelID;
            PanelTitle = panelTitle;
            PanelContent = panelContent;
            PanelRecommendedWidth = panelRecommendedWidth;
            PanelRecommendedHeight = panelRecommendedHeight;
            PanelMinWidth = panelMinWidth;
            PanelMinHeight = panelMinHeight;
            PanelMaxWidth = panelMaxWidth;
            PanelMaxHeight = panelMaxHeight;
            PersistentPanelWidth = persistentPanelWidth;
            PersistentPanelHeight = persistentPanelHeight;
            PanelRecommendedX = panelRecommendedX;
            PanelRecommendedY = panelRecommendedY;
            PersistentPanelX = persistentPanelX;
            PersistentPanelY = persistentPanelY;
            ControlPort = controlPort;
            Deleted = deleted;
            PersistentPanelData = persistentPanelData;
            PanelIcon = panelIcon;
            ClientID = clientID;
        }

        public static RuntimeCustomPanel FromJsonElement(JsonElement x)
        {
            RuntimeCustomPanel runtimeCustomPanel = new RuntimeCustomPanel(
        x.GetProperty("PanelType").GetString()!,
    x.GetProperty("LoaderPanelID").GetString()!,
        x.GetProperty("PersistentPanelID").GetString()!,
   x.GetProperty("PanelTitle").GetString()!,
   x.GetProperty("PanelContent").GetString()!,
   x.GetProperty("PanelRecommendedWidth").GetDouble(),
   x.GetProperty("PanelRecommendedHeight").GetDouble(),
   x.GetProperty("PanelMinWidth").GetDouble(),
   x.GetProperty("PanelMinHeight").GetDouble(),
   x.GetProperty("PanelMaxWidth").GetDouble(),
   x.GetProperty("PanelMaxHeight").GetDouble(),
   x.GetProperty("PersistentPanelWidth").GetDouble(),
   x.GetProperty("PersistentPanelHeight").GetDouble(),
      x.GetProperty("PanelRecommendedX").GetDouble(),
   x.GetProperty("PanelRecommendedY").GetDouble(),
   x.GetProperty("PersistentPanelX").GetDouble(),
   x.GetProperty("PersistentPanelY").GetDouble(),
   x.GetProperty("ControlPort").GetInt32(),
   x.GetProperty("Deleted").GetBoolean(),
   x.GetProperty("PersistentPanelData").GetString()!,
   x.GetProperty("PanelIcon").GetString()!,
   x.GetProperty("ClientID").GetString()!
    );
            return runtimeCustomPanel;
        }

        public RuntimeCustomPanel Cloned()
        {
            var serialized = JsonSerializer.Serialize( this );
            return JsonSerializer.Deserialize<RuntimeCustomPanel>( serialized );
        }
    }

    public class TemplateCustomPanel
    {
        public string PanelType { get; }
        public string LoaderPanelID { get; }
        public string PanelTitle { get; }
        public string PanelContent { get; }
        public double PanelRecommendedWidth { get; }
        public double PanelRecommendedHeight { get; }
        public double PanelMinWidth { get; }
        public double PanelMinHeight { get; }
        public double PanelMaxWidth { get; }
        public double PanelMaxHeight { get; }
        public double PanelRecommendedX { get; }
        public double PanelRecommendedY { get; }
        public string PanelDefaultData { get; }
        public string PanelIcon { get; }
        public string ClientID { get; }

        public TemplateCustomPanel(string panelType, string loaderPanelID, string panelTitle, string panelContent, double panelRecommendedWidth, double panelRecommendedHeight, double panelMinWidth, double panelMinHeight, double panelMaxWidth, double panelMaxHeight, double panelRecommendedX, double panelRecommendedY, string panelDefaultData, string panelIcon, string clientID)
        {
            PanelType = panelType;
            LoaderPanelID = loaderPanelID;
            PanelTitle = panelTitle;
            PanelContent = panelContent;
            PanelRecommendedWidth = panelRecommendedWidth;
            PanelRecommendedHeight = panelRecommendedHeight;
            PanelMinWidth = panelMinWidth;
            PanelMinHeight = panelMinHeight;
            PanelMaxWidth = panelMaxWidth;
            PanelMaxHeight = panelMaxHeight;
            PanelRecommendedX = panelRecommendedX;
            PanelRecommendedY = panelRecommendedY;
            PanelDefaultData = panelDefaultData;
            PanelIcon = panelIcon;
            ClientID = clientID;
        }

        public static TemplateCustomPanel FromJsonElement(JsonElement x)
        {
            TemplateCustomPanel templateCustomPanel = new TemplateCustomPanel(
                    x.GetProperty("PanelType").GetString()!,
                    x.GetProperty("LoaderPanelID").GetString()!,
                    x.GetProperty("PanelTitle").GetString()!,
                    x.GetProperty("PanelContent").GetString()!,
                    x.GetProperty("PanelRecommendedWidth").GetDouble(),
                    x.GetProperty("PanelRecommendedHeight").GetDouble(),
                    x.GetProperty("PanelMinWidth").GetDouble(),
                    x.GetProperty("PanelMinHeight").GetDouble(),
                    x.GetProperty("PanelMaxWidth").GetDouble(),
                    x.GetProperty("PanelMaxHeight").GetDouble(),
                    x.GetProperty("PanelRecommendedX").GetDouble(),
                    x.GetProperty("PanelRecommendedY").GetDouble(),
                    x.GetProperty("PanelDefaultData").GetString()!,
                    x.GetProperty("PanelIcon").GetString()!,
                    x.GetProperty("ClientID").GetString()!
                );
            return templateCustomPanel;
        }

        public TemplateCustomPanel Cloned()
        {
            var serialized = JsonSerializer.Serialize( this );
            return JsonSerializer.Deserialize<TemplateCustomPanel>(serialized);
        }
    }

    public class PanelHeader
    {
        public string persistentpanelid { get; set; }
        public bool titlebarvisible { get; set; }
        public bool enableclose { get; set; }
        public bool enableresize { get; set; }
        public bool enabledrag { get; set; }
    }
}
