class OutlyerCli < Formula
  desc "Outlyer CLI allows to easily manage your Outlyer account via command line."
  homepage "https://www.outlyer.com/"
  url "https://github.com/outlyerapp/outlyer-cli/releases/download/1.0.1/outlyer-cli_1.0.1_Darwin_x86_64.tar.gz"
  version "1.0.1"
  sha256 "18c7263a27975572ee27978fff501bb5b85c7718b7bc14fa42cf24697ef72153"

  def install
    bin.install "outlyer"
  end
end
